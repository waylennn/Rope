package registry

import (
	"context"
	"fmt"
	"sync"
)

//PluginMgr ...
type PluginMgr struct {
	lock    sync.Mutex
	plugins map[string]Registry
}

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
<<<<<<< HEAD
=======

>>>>>>> 4f008cdac20aeb2e116d997fc89c3568a6dca67a
	p.lock.Lock()
	defer p.lock.Unlock()
	_, ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("duplicate registry plugin")
		return
	}

	p.plugins[plugin.Name()] = plugin
<<<<<<< HEAD
=======

>>>>>>> 4f008cdac20aeb2e116d997fc89c3568a6dca67a
	return
}

func (p *PluginMgr) initRegistry(ctx context.Context, name string, opt ...Option) (registry Registry, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	plugin, ok := p.plugins[name]
	if !ok {
		err = fmt.Errorf("plugin is not exist")
		return
	}

	registry = plugin
	err = registry.Init(ctx, opt...)
	return
}

//RegisterPlugin external interface
func RegisterPlugin(plugin Registry) (err error) {
	return pluginMgr.registerPlugin(plugin)
}

//InitRegistry external interface
func InitRegistry(ctx context.Context, name string, opt ...Option) (registry Registry, err error) {
	return pluginMgr.initRegistry(ctx, name, opt...)
}
