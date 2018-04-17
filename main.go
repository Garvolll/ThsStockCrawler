package main

import (
	"fmt"
	"sourcegraph.com/sourcegraph/go-selenium"
	"time"
	"github.com/Luxurioust/excelize"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type stockItem struct {
	Code               string
	Name               string
	RiseFallRate       string
	ChangeRate         string
	LastestPrice       string
	PreviousPoints     string
	PreviousPointsDate string
	CreateDate         string
	DataTypes          int
}

const GURL = "http://data.10jqka.com.cn/rank/cxg/"
const DURL = "http://data.10jqka.com.cn/rank/cxd/"

var webDriver selenium.WebDriver
var err error
var xlsx = excelize.NewFile()
var db *sql.DB

func GetBySelenium(url string) {

	err = webDriver.Get(url)
	if err != nil {
		fmt.Printf("Failed to load page: %s\n", err)
		return
	}
	elem, err := webDriver.FindElement(selenium.ByCSSSelector, "#datacenter_change_content > div.table-tab.J-ajax-board > a:nth-child(3)")
	if err != nil {
		fmt.Printf("Failed to find button : %s\n", err)
		return
	}
	elem.Click()
	time.Sleep(1000)
	trs, err := webDriver.FindElements(selenium.ByXPATH, `//*[@id="J-ajax-main"]/table/tbody/tr`)
	if err != nil {
		fmt.Printf("Failed to find data item: %s\n", err)
		return
	}
	fmt.Printf("共有%d只股票\n", len(trs))

	var dataType int
	if url == "http://data.10jqka.com.cn/rank/cxg/" {
		dataType = 1
	} else {
		dataType = 2
	}

	tx, _ := db.Begin()
	//每次循环用的都是tx内部的连接，没有新建连接，效率高
	tx.Exec("INSERT INTO gd_count(date,count,data_type)"+
		"values(?,?,?)", time.Now().Format("2006-01-02"), len(trs), dataType)
	//最后释放tx内部的连接
	tx.Commit()


	//for i := 0; i < len(trs); i++ {
	//	//go HandleItemConcurrent(resChan)
	//	handleItem(dataType, trs[i], i)
	//}
	//date := time.Now().String()
	//// Save xlsx file by the given path.
	//xerr := xlsx.SaveAs("./" + date[0:10] + ".xlsx")
	//if xerr != nil {
	//	fmt.Println(err)
	//}

}

func handleItem(dataType string, item selenium.WebElement, index int) {
	tds, _ := item.FindElements(selenium.ByXPATH, "td")
	s := stockItem{}
	codeTd, _ := tds[1].FindElement(selenium.ByXPATH, "a")
	s.Code, _ = codeTd.Text()
	nameTd, _ := tds[2].FindElement(selenium.ByXPATH, "a")
	s.Name, _ = nameTd.Text()
	s.RiseFallRate, _ = tds[3].Text()
	s.ChangeRate, _ = tds[4].Text()
	s.LastestPrice, _ = tds[5].Text()
	s.PreviousPoints, _ = tds[6].Text()
	PreviousPointsDate, _ := tds[7].Text()
	/**chusheng**/
	timeLayout := "2006-01-02"           //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, PreviousPointsDate, loc)
	//timestamp := tmp.Unix()    //转化为时间戳 类型是int64
	s.PreviousPointsDate = tmp.Format("2006-01-02")
	/****/
	s.CreateDate = time.Now().Format("2006-01-02")
	fmt.Println(s)
	//var sheet string
	//if dataType == "G" {
	//	sheet = "Sheet1"
	//	s.DataTypes = 1
	//} else {
	//	sheet = "Sheet2"
	//	s.DataTypes = 2
	//	xlsx.NewSheet("Sheet2")
	//}

	tx, _ := db.Begin()
	//每次循环用的都是tx内部的连接，没有新建连接，效率高
	tx.Exec("INSERT INTO gd(code,s_name,risefall_rate,change_rate,lastest_price,previous_points,previous_points_date,create_time,data_types)"+
		"values(?,?,?,?,?,?,?,?,?)", s.Code, s.Name, s.RiseFallRate, s.ChangeRate, s.LastestPrice, s.PreviousPoints, s.PreviousPointsDate, s.CreateDate, s.DataTypes)
	//最后释放tx内部的连接
	tx.Commit()

	//_index := strconv.Itoa(index + 2)
	//xlsx.SetCellValue(sheet, "A1", "股票代码")
	//xlsx.SetCellValue(sheet, "B1", "股票简称")
	//xlsx.SetCellValue(sheet, "C1", "涨跌幅")
	//xlsx.SetCellValue(sheet, "D1", "换手率")
	//xlsx.SetCellValue(sheet, "E1", "最新价(元)")
	//xlsx.SetCellValue(sheet, "F1", "前期高点")
	//xlsx.SetCellValue(sheet, "G1", "前期高点日期")
	//
	//xlsx.SetCellValue(sheet, "A"+_index, s.Code)
	//xlsx.SetCellValue(sheet, "B"+_index, s.Name)
	//xlsx.SetCellValue(sheet, "C"+_index, s.RiseFallRate)
	//xlsx.SetCellValue(sheet, "D"+_index, s.ChangeRate)
	//xlsx.SetCellValue(sheet, "E"+_index, s.LastestPrice)
	//xlsx.SetCellValue(sheet, "F"+_index, s.PreviousPoints)
	//xlsx.SetCellValue(sheet, "G"+_index, s.PreviousPointsDate)

}

func main() {
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox"})
	//配置chromeDriver所使用的端口，默认 http://localhost:9515
	if webDriver, err = selenium.NewRemote(caps, "http://localhost:9515"); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/stock")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	defer webDriver.Quit()

	GetBySelenium(GURL)
	GetBySelenium(DURL)

}
