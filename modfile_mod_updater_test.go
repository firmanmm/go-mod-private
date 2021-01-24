package gomodprivate

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModfileModUpdater_Update(t *testing.T) {
	type args struct {
		repositories []string
		inputContent string
	}

	type want struct {
		name    string
		version string
	}

	tests := []struct {
		name    string
		args    args
		wants   []want
		wantErr bool
	}{
		{
			"ValidFlow",
			args{
				[]string{
					"rendoru.com/module/sync-mq@v1.0.0",
					"rendoru.com/tool/concurrent",
				},
				`
				module rendoru.com/service/dorurumon

				go 1.12
				
				require (
					github.com/Masterminds/squirrel v1.4.0
					github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
					github.com/firmanmm/gin-easy-route v0.0.0-20190717143707-8951ff5013e7
					github.com/firmanmm/go-mirror v1.0.1
					github.com/firmanmm/go-templater v1.0.0
					github.com/firmanmm/suberror v1.2.0
					github.com/gin-gonic/gin v1.6.3
					github.com/go-ini/ini v1.62.0
					github.com/go-playground/validator/v10 v10.4.1 // indirect
					github.com/go-redis/redis v6.15.9+incompatible
					github.com/go-sql-driver/mysql v1.5.0
					github.com/golang/protobuf v1.4.3 // indirect
					github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
					github.com/json-iterator/go v1.1.10
					github.com/onsi/ginkgo v1.14.2 // indirect
					github.com/onsi/gomega v1.10.3 // indirect
					github.com/russross/blackfriday/v2 v2.1.0 // indirect
					github.com/sarulabs/di v2.0.0+incompatible
					github.com/smartystreets/assertions v1.2.0 // indirect
					github.com/smartystreets/goconvey v1.6.4 // indirect
					github.com/streadway/amqp v1.0.0
					github.com/ugorji/go v1.2.0 // indirect
					github.com/urfave/cli v1.22.5
					go.uber.org/fx v1.13.1
					golang.org/x/crypto v0.0.0-20201112155050-0c6587e931a9 // indirect
					golang.org/x/sys v0.0.0-20201113233024-12cec1faf1ba // indirect
					golang.org/x/tools v0.0.0-20201113202037-1643af1435f3 // indirect
					google.golang.org/protobuf v1.25.0 // indirect
					gopkg.in/ini.v1 v1.62.0 // indirect
				)			
				`,
			},
			[]want{
				{
					"rendoru.com/module/sync-mq",
					"v1.0.0",
				},
				{
					"rendoru.com/tool/concurrent",
					"v0.0.0",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewModfileModUpdater()
			testFileName := "TestModfileModUpdater_Update.go.mod.tf"
			g.SetTargetFile(testFileName)
			err := ioutil.WriteFile(testFileName, []byte(tt.args.inputContent), 0666)
			assert.NoError(t, err)
			defer os.Remove(testFileName)
			if err := g.Update(tt.args.repositories); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			rawFinalRead, err := ioutil.ReadFile(testFileName)
			assert.NoError(t, err)
			finalRead := string(rawFinalRead)
			fmt.Println(finalRead)
			for _, want := range tt.wants {
				assert.Contains(t, finalRead, want.name)
				assert.Contains(t, finalRead, fmt.Sprintf("%s %s => ./.vendor.gomp/%s", want.name, want.version, want.name))
			}
		})
	}
}
