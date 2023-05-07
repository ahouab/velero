package schedule

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/vmware-tanzu/velero/test/e2e"
	. "github.com/vmware-tanzu/velero/test/e2e/test"
	. "github.com/vmware-tanzu/velero/test/e2e/util/k8s"
	. "github.com/vmware-tanzu/velero/test/e2e/util/velero"
)

type ScheduleBackupCreation struct {
	TestCase
	namespace        string
	ScheduleName     string
	ScheduleArgs     []string
	Period           int //Limitation: The unit is minitue only and 60 is divisible by it
	randBackupName   string
	verifyTimes      int
	volume           string
	podName          string
	pvcName          string
	podAnn           map[string]string
	podSleepDuration time.Duration
}

var ScheduleBackupCreationTest func() = TestFunc(&ScheduleBackupCreation{})

func (n *ScheduleBackupCreation) Init() error {
	n.TestCase.Init()
	n.VeleroCfg = VeleroCfg
	n.Client = *n.VeleroCfg.ClientToInstallVelero
	n.Period = 3      // Unit is minute
	n.verifyTimes = 5 // More larger verify times more confidence we have"
	podSleepDurationStr := "300s"
	n.podSleepDuration, _ = time.ParseDuration(podSleepDurationStr)
	n.CaseBaseName = "schedule-backup-creation-test" + n.UUIDgen
	n.namespace = n.CaseBaseName

	n.TestMsg = &TestMSG{
		Desc:      "Schedule controller wouldn't create a new backup when it still has pending or InProgress backup",
		FailedMSG: "Failed to verify schedule back creation behavior",
		Text:      "Schedule controller wouldn't create a new backup when it still has pending or InProgress backup",
	}
	n.podAnn = map[string]string{
		"pre.hook.backup.velero.io/container": n.podName,
		"pre.hook.backup.velero.io/command":   "[\"sleep\", \"" + podSleepDurationStr + "\"]",
		"pre.hook.backup.velero.io/timeout":   "600s",
	}
	n.volume = "volume-1"
	n.podName = "pod-1"
	n.pvcName = "pvc-1"
	n.ScheduleName = "schedule-" + n.CaseBaseName

	n.ScheduleArgs = []string{
		"--include-namespaces", n.namespace,
		"--schedule=*/" + fmt.Sprintf("%v", n.Period) + " * * * *",
	}
	Expect(n.Period < 30).To(Equal(true))
	return nil
}

func (s *ScheduleBackupCreation) CreateResources() error {
	s.Ctx, s.CtxCancel = context.WithTimeout(context.Background(), 60*time.Minute)
	By(fmt.Sprintf("Create namespace %s", s.namespace), func() {
		Expect(CreateNamespace(s.Ctx, s.Client, s.namespace)).To(Succeed(),
			fmt.Sprintf("Failed to create namespace %s", s.namespace))
	})

	By(fmt.Sprintf("Create pod %s in namespace %s", s.podName, s.namespace), func() {
		_, err := CreatePod(s.Client, s.namespace, s.podName, "default", s.pvcName, []string{s.volume}, nil, s.podAnn)
		Expect(err).To(Succeed())
		err = WaitForPods(s.Ctx, s.Client, s.namespace, []string{s.podName})
		Expect(err).To(Succeed())
	})
	return nil
}

func (n *ScheduleBackupCreation) Backup() error {
	// Wait until the beginning of the given period to create schedule, it will give us
	//   a predictable period to wait for the first scheduled backup, and verify no immediate
	//   scheduled backup was created between schedule creation and first scheduled backup.
	By(fmt.Sprintf("Creating schedule %s ......\n", n.ScheduleName), func() {
		for i := 0; i < n.Period*60/30; i++ {
			time.Sleep(30 * time.Second)
			now := time.Now().Minute()
			triggerNow := now % n.Period
			if triggerNow == 0 {
				Expect(VeleroScheduleCreate(n.Ctx, VeleroCfg.VeleroCLI, VeleroCfg.VeleroNamespace, n.ScheduleName, n.ScheduleArgs)).To(Succeed(), func() string {
					RunDebug(context.Background(), VeleroCfg.VeleroCLI, VeleroCfg.VeleroNamespace, "", "")
					return "Fail to create schedule"
				})
				break
			}
		}
	})

	By("Delay one more minute to make sure the new backup was created in the given period", func() {
		time.Sleep(1 * time.Minute)
	})

	By(fmt.Sprintf("Get backups every %d minute, and backups count should increase 1 more step in the same pace\n", n.Period), func() {
		for i := 1; i <= n.verifyTimes; i++ {
			fmt.Printf("Start to sleep %f minute #%d time...\n", n.podSleepDuration.Minutes(), i)
			mi, _ := time.ParseDuration("60s")
			time.Sleep(n.podSleepDuration + mi)
			bMap := make(map[string]string)
			backupsInfo, err := GetScheduledBackupsCreationTime(n.Ctx, VeleroCfg.VeleroCLI, "default", n.ScheduleName)
			Expect(err).To(Succeed())
			Expect(len(backupsInfo) == i).To(Equal(true))
			for index, bi := range backupsInfo {
				bList := strings.Split(bi, ",")
				fmt.Printf("Backup %d: %v\n", index, bList)
				bMap[bList[0]] = bList[1]
				_, err := time.Parse("2006-01-02 15:04:05 -0700 MST", bList[1])
				Expect(err).To(Succeed())
			}
			if i == n.verifyTimes-1 {
				backupInfo := backupsInfo[rand.Intn(len(backupsInfo))]
				n.randBackupName = strings.Split(backupInfo, ",")[0]
			}
		}
	})
	return nil
}

func (n *ScheduleBackupCreation) Clean() error {
	if !n.VeleroCfg.Debug {
		Expect(VeleroScheduleDelete(n.Ctx, n.VeleroCfg.VeleroCLI, n.VeleroCfg.VeleroNamespace, n.ScheduleName)).To(Succeed())
		Expect(n.TestCase.Clean()).To(Succeed())
	}
	return nil
}
