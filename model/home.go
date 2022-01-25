package model

type BasicMenu struct {
	Id int
	AuthName string
	Path string
}

func NewBasicMenu(newId int,newAuthName string,newPath string) BasicMenu {
	return BasicMenu{
		Id:       newId,
		AuthName: newAuthName,
		Path:     newPath,
	}
}

type MenuListItem struct {
	Father BasicMenu
	Children []BasicMenu
}

func NewMenuListItem(id int,name string, path string,childrenLength int) MenuListItem {
	return MenuListItem{
		Father:       NewBasicMenu(id,name,path),
		// 动态分配子数组长度
		Children: make([]BasicMenu,childrenLength),
	}
}

type MenuList struct {
	Menu  []MenuListItem
}



func NewMenu() []MenuListItem {
	//MenuList:=NewMenuList()
	m:=MenuList{Menu:make([]MenuListItem,5,10)}

	m.Menu[0] = NewMenuListItem(100,"用户管理","user",1)
	m.Menu[0].Children[0] = NewBasicMenu(101,"用户列表","users")

	m.Menu[1] = NewMenuListItem(110,"权限管理","admin",2)
	m.Menu[1].Children[0] = NewBasicMenu(111,"角色列表","roles")
	m.Menu[1].Children[1] = NewBasicMenu(112,"权限列表","rights")


	m.Menu[2] = NewMenuListItem(120,"商品管理","goods",1)
	m.Menu[2].Children[0] = NewBasicMenu(121,"商品详情","goods")

	m.Menu[3] = NewMenuListItem(130,"订单管理","orders",2)
	m.Menu[3].Children[0] = NewBasicMenu(131,"订单增删改查","ordersCURD")
	m.Menu[3].Children[1] = NewBasicMenu(132,"点单数","orderNums")

	m.Menu[4]= NewMenuListItem(140,"数据统计","data",2)
	m.Menu[4].Children[0] = NewBasicMenu(142,"折线图","gic")
	m.Menu[4].Children[1] = NewBasicMenu(143,"汇总","sum")


	return m.Menu

}