package main

//AllDirList 需要创建的目录
// var AllDirList []string = []string{
// 	"controller",
// 	"idl",
// 	"main",
// 	"scripts",
// 	"conf",
// 	"app/router",
// 	"app/config",
// 	"model",
// 	"generate",
// }

//DirGenerator 创建目录
type DirGenerator struct {
	dirList []string
}

//Run .
func (d *DirGenerator) Run(opt *Option, metaData *ServiceMetaData) error {

	// for _, dir := range d.dirList {
	// 	fullDir := path.Join(opt.Output, dir)
	// 	err := os.MkdirAll(fullDir, 0755)
	// 	if err != nil {
	// 		fmt.Printf("mkdir dir %s failed, err:%v\n", dir, err)
	// 		return err
	// 	}
	// }

	return nil
}

func init() {
	// dir := &DirGenerator{
	// 	dirList: AllDirList,
	// }

	// err := Register("dir generator", dir)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
