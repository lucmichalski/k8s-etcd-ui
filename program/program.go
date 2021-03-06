package program

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/etcd-manage/etcd-manage-server/program/api"
	v1 "github.com/etcd-manage/etcd-manage-server/program/api/v1"
	"github.com/etcd-manage/etcd-manage-server/program/config"
	"github.com/etcd-manage/etcd-manage-server/program/logger"
	"github.com/etcd-manage/etcd-manage-server/program/models"
)

// Program 主程序
type Program struct {
	cfg   *config.Config     // 配置
	s     *http.Server       // http服务-用于程序结束stop
	vApis map[string]api.API // 多版本api启动运行
}

// New 创建主程序
func New() (*Program, error) {
	// 配置文件
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 日志对象
	_, err = logger.InitLogger(cfg.LogPath, cfg.Debug)
	if err != nil {
		return nil, err
	}

	// 连接数据库
	err = models.InitClient(cfg.DB)
	if err != nil {
		return nil, err
	}

	// 注册api对象
	apis := make(map[string]api.API, 0)
	apis["v1"] = new(v1.APIV1)

	// js, _ := json.Marshal(cfg)
	// fmt.Println(string(js))

	return &Program{
		cfg:   cfg,
		vApis: apis,
	}, nil
}

// Run 启动程序
func (p *Program) Run() error {
	// 启动http服务
	go p.startAPI()

	// 打开浏览器
	go func() {
		time.Sleep(100 * time.Millisecond)
		openUrl(fmt.Sprintf("http://127.0.0.1:%d/ui/", p.cfg.HTTP.Port))
	}()

	return nil
}

// Stop 停止服务
func (p *Program) Stop() {
	if p.s != nil {
		p.s.Close()
	}
}

func openUrl(uri string) error {
	var (
		commands = map[string]string{
			"windows": "start",
			"darwin":  "open",
			"linux":   "xdg-open",
		}
	)
	run, ok := commands[runtime.GOOS] //获取平台信息

	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	cmd := exec.Command(run, uri)
	return cmd.Start()
}
