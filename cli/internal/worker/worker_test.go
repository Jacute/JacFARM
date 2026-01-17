package worker

import (
	jacfarm_client "cli_exploit_runner/internal/clients/jacfarm"
	"cli_exploit_runner/internal/worker/mocks"
	"context"
	"log/slog"
	"net"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/jacute/prettylogger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAttack(t *testing.T) {
	const testPath = "testcases/testsploit.sh"

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Fatal(err)
	}

	if err := os.Chmod(testPath, 0744); err != nil {
		t.Fatal(err)
	}

	out, err := attack(context.Background(), testPath, "127.0.0.1")
	require.NoError(t, err)

	flags := parseFlags(out, regexp.MustCompile("[A-Z0-9]{31}="))
	require.Equal(t, []string{
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC=",
	}, flags)
}

func TestAttackAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testcases := []struct {
		name   string
		flags  []string
		worker func() *Worker
	}{
		{
			name: "ok",
			flags: []string{
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC=",
			},
			worker: func() *Worker {
				clientMock := mocks.NewMockJacFARMClient(ctrl)
				clientMock.EXPECT().GetTeams(gomock.Any()).Return([]*jacfarm_client.Team{
					{
						ID:   1,
						Name: "aboba",
						IP:   net.ParseIP("1.1.1.1"),
					},
					{
						ID:   2,
						Name: "aboba2",
						IP:   net.ParseIP("1.1.1.2"),
					},
				}, nil).Times(2)
				clientMock.EXPECT().SendFlags(gomock.Any(), mocks.UnorderedSlice([]*jacfarm_client.ServiceFlag{
					{
						Flag:   "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
						TeamID: 1,
					},
					{
						Flag:   "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=",
						TeamID: 1,
					},
					{
						Flag:   "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC=",
						TeamID: 1,
					},
					{
						Flag:   "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
						TeamID: 2,
					},
					{
						Flag:   "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=",
						TeamID: 2,
					},
					{
						Flag:   "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC=",
						TeamID: 2,
					},
				})).Return(nil).Times(1)
				w, err := New(
					clientMock,
					slog.New(prettylogger.NewDiscardHandler()),
					"testcases/testsploit.sh",
					"[A-Z0-9]{31}=",
					WithAttackPeriod(200*time.Millisecond),
				)
				if err != nil {
					t.Fatal(err)
				}

				return w
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(tt *testing.T) {
			w := tc.worker()
			w.Run()

			time.Sleep(1 * time.Second)
			w.Stop()
		})
	}
}
