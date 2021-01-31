package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"github.com/idoubi/goz"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)
type Option struct {
	Id string
	Title string
}
type Answer struct {
	Activity_id string `json:"activity_id"`
	Question_id string `json:"question_id"`
	Answer []string	`json:"answer"`
	Mode_id string	`json:"mode_id"`
	Way string	`json:"way"`
}
type Options struct {
	Count int
	Option []Option
}
type Finish struct {
	Race_code string `json:"race_code"`
}
type People struct {
	name string
	province string
	univ string
	correct int64
	consume int64
}
var (
	kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)

	// 给字体颜色对象赋值
	FontColor Color = Color{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
)

type Color struct {
	black        int // 黑色
	blue         int // 蓝色
	green        int // 绿色
	cyan         int // 青色
	red          int // 红色
	purple       int // 紫色
	yellow       int // 黄色
	light_gray   int // 淡灰色（系统默认值）
	gray         int // 灰色
	light_blue   int // 亮蓝色
	light_green  int // 亮绿色
	light_cyan   int // 亮青色
	light_red    int // 亮红色
	light_purple int // 亮紫色
	light_yellow int // 亮黄色
	white        int // 白色
}

// 输出有颜色的字体
func ColorPrint(s string, i int) {
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	print(s)
	CloseHandle.Call(handle)
}

func toChinese(textUnquoted string ) string {
	sUnicodev := strings.Split(textUnquoted, "\\u")
	var context string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)
	}
	return(context)
}
func _auth() string {
	file, err := os.Open("./auth.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	return string(content)
}
func main() {
	var auth,act_id string
	auth = _auth()
	act_id = "5f71e934bcdbf3a8c3ba5061"
	cli := goz.NewClient()

	resp, err := cli.Get("https://ssxx.univs.cn/cgi-bin/race/grade/?t="+strconv.FormatInt(int64(time.Now().Local().Unix()),10)+"&activity_id="+act_id, goz.Options{
		Headers: map[string]interface{}{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
			"Authorization": auth,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := resp.GetBody()
	name , jerr := jsonparser.GetString(body,"data","name")
	province , _  := jsonparser.GetString(body,"data","province_name")
	university , _ := jsonparser.GetString(body,"data","university_name")
	integral , _  := jsonparser.GetInt(body,"data","integral")
	join_times , _  := jsonparser.GetInt(body,"data","join_times")
	if  jerr == nil {
		ColorPrint("[info] 成功获取登录信息\n",FontColor.light_blue);
		ColorPrint("[info] 姓名："+name+"\n",FontColor.light_blue);
		ColorPrint("[info] 学校：["+province+"]"+university+"\n",FontColor.light_blue);
		ColorPrint("[info] 已获得积分："+strconv.FormatInt(int64(integral),10)+"\n",FontColor.light_blue)
		ColorPrint("[info] 已比赛次数："+strconv.FormatInt(int64(join_times),10)+"\n",FontColor.light_blue)
	} else {
		ColorPrint("[error] 您的 Authorization Token 不正确或已过期\n",FontColor.light_red)
		fmt.Println("按任意键退出...")
		var input string
		fmt.Scanln(&input)
		return
	}
	mode := make(map[int]map[int]string)
	for i:=1;i<=4;i++{
		mode[i] = make(map[int]string)
	}
	mode[1][0]="[1]限时赛";mode[1][1]="5f71e934bcdbf3a8c3ba51d9"
	mode[2][0]="[2]抢十赛";mode[2][1]="5f71e934bcdbf3a8c3ba51da"
	ColorPrint("题目类型\n",FontColor.light_gray)
	for i:=1;i<=2;i++{
		ColorPrint(mode[i][0]+"\n",FontColor.light_gray)
	}
	ColorPrint("请输入您比赛类型序号：",FontColor.light_gray)
	var mid int
	fmt.Scanln(&mid)
	_mid := mode[mid][1]
	//println(_mid)
	//begin
	var i1,times,i2,wait int
	ColorPrint("请输入循环比赛次数：",FontColor.light_gray)
	fmt.Scanln(&times)
	ColorPrint("请输入答题间隔：",FontColor.light_gray)
	fmt.Scanln(&wait)
	for i1=1;i1<=times;i1++ {
		i2=0
		ColorPrint("[info] 现在开始第"+strconv.FormatInt(int64(i1),10)+"轮比赛\n",FontColor.light_blue)
		url := "https://ssxx.univs.cn/cgi-bin/race/beginning/?t=" + strconv.FormatInt(int64(time.Now().Local().Unix()), 10) + "&activity_id=" + act_id + "&mode_id=" + _mid + "&way=1"
		resp, err = cli.Get(url, goz.Options{
			Headers: map[string]interface{}{
				"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
				"Authorization": auth,
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
		body, _ = resp.GetBody()
		rid, _ := jsonparser.GetString(body, "race_code")
        var myself,opponent People
        myself.name,_ = jsonparser.GetString(body,"myself","name")
		myself.province,_ = jsonparser.GetString(body,"myself","province_name")
		myself.univ,_ = jsonparser.GetString(body,"myself","univ_name")
		opponent.name,_ = jsonparser.GetString(body,"opponent","name")
		opponent.province,_ = jsonparser.GetString(body,"opponent","province_name")
		opponent.univ,_ = jsonparser.GetString(body,"opponent","univ_name")
		ColorPrint("[info] 我方选手["+myself.province+"]"+myself.univ+":"+myself.name+"\n",FontColor.light_blue)
		ColorPrint("[info] 对方选手["+opponent.province+"]"+opponent.univ+":"+opponent.name+"\n",FontColor.light_blue)
		jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			i2++
			qid := string(value)
			//fmt.Println(qid)
			url = "https://ssxx.univs.cn/cgi-bin/race/question/?t=" + strconv.FormatInt(int64(time.Now().Local().Unix()), 10) + "&activity_id=" + act_id + "&question_id=" + qid + "&mode_id=" + _mid + "&way=1"
			resp, err = cli.Get(url, goz.Options{
				Headers: map[string]interface{}{
					"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
					"Authorization": auth,
				},
			})
			if err != nil {
				log.Fatalln(err)
			}
			ColorPrint("[info] 第"+strconv.FormatInt(int64(i2),10)+"题",FontColor.light_blue)
			content, _ := resp.GetBody()
			title, _ := jsonparser.GetString(content, "data", "title")
			id, _ := jsonparser.GetString(content, "data", "id")
			//fmt.Println(string(content))
			db, err := sql.Open("sqlite3", "./question.db")
			var count, count1 int64
			var SAns Answer
			db.QueryRow("select count(*) from question where `id` = \"" + id + "\"").Scan(&count)
			db.QueryRow("select count(*) from question where `title` = \"" + getVisable(title) + "\"").Scan(&count1)
			if (count == 0) {
				stmt, _ := db.Prepare("INSERT INTO `question`(`id`, `title`, `answer`,`options`)  values(?, ?, ?,?)")
				stmt.Exec(id, getVisable(title), "", parseAns(content))
			}
			if (count1 != 0) { //记得修改
				ColorPrint(",已找到答案",FontColor.light_green)
				stmt, _ := db.Prepare("UPDATE `question` SET `options` = ? WHERE `id` =?")
				stmt.Exec(parseAns(content), id)
				var ansmap map[string]string
				ansmap = make(map[string]string)
				rows, _ := db.Query("SELECT `answer` FROM `question` where `title`=\"" + getVisable(title) + "\"")
				for rows.Next() {
					var ans string
					rows.Scan(&ans)
					jsonparser.ArrayEach([]byte(ans), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						tit, _ := jsonparser.GetString(value, "Title")
						ansmap[_md5(tit)] = tit
					}, "Option")
				}
				//fmt.Println(ansmap)
				SAns = findAns(content, ansmap, act_id, _mid, id)
				if (SAns.Answer == nil) {
					ColorPrint(",数据库中答案异常，重设答案",FontColor.light_yellow)
					tmp, _ := jsonparser.GetString(content, "data", "options", "[0]", "id")
					SAns.Answer = append(SAns.Answer, tmp)
				}
			} else {
				var Ans Answer
				Ans.Activity_id = act_id
				Ans.Way = "1"
				Ans.Mode_id = _mid
				Ans.Question_id = id
				tmp, _ := jsonparser.GetString(content, "data", "options", "[0]", "id")
				Ans.Answer = append(Ans.Answer, tmp)
				SAns = Ans
				ColorPrint(",未找到答案",FontColor.light_yellow)
			}
			//bytes, _ := json.Marshal(SAns)
			//fmt.Println(string(bytes))
			url = "https://ssxx.univs.cn/cgi-bin/race/answer/"
			resp, err = cli.Post(url, goz.Options{
				Headers: map[string]interface{}{
					"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
					"Authorization": auth,
					"Content-Type":  "application/json;charset=UTF-8",
				},
				JSON: Answer(SAns),
			})
			if err != nil {
				log.Fatalln(err)
			}
			cans, _ := resp.GetBody()
			//fmt.Println(cans)
			getAns(cans, id)
			correct, _ := jsonparser.GetBoolean(cans, "data", "correct")
			//fmt.Print(correct)
			if correct == true {
				ColorPrint(",回答正确\n", FontColor.light_green)
			} else {
				ColorPrint(",回答错误\n", FontColor.light_red)
			}
			time.Sleep(time.Duration(wait) * time.Second)
		}, "question_ids")

		var _Finish Finish
		_Finish.Race_code = rid
		url = "https://ssxx.univs.cn/cgi-bin/race/finish/"
		resp, err = cli.Post(url, goz.Options{
			Headers: map[string]interface{}{
				"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
				"Authorization": auth,
				"Content-Type":  "application/json;charset=UTF-8",
			},
			JSON: Finish(_Finish),
		})
		finish, _ := resp.GetBody()
		//fmt.Println(string(finish))
		ccount, _ := jsonparser.GetInt(finish, "data", "owner", "correct_amount")
		myself.correct=int64(ccount)
		ccount1, _ := jsonparser.GetInt(finish, "data", "opponent", "correct_amount")
		opponent.correct=int64(ccount1)
		myself.consume ,_ =jsonparser.GetInt(finish,"data","owner","consume_time")
		opponent.consume ,_ =jsonparser.GetInt(finish,"data","opponent","consume_time")
		ColorPrint("[info] 答题结束，我方用时"+strconv.FormatInt(int64(myself.consume), 10)+"秒,共回答正确：" + strconv.FormatInt(int64(myself.correct), 10) + "题，正确率",FontColor.light_gray)
		fmt.Printf("%.2f",float64(myself.correct)/float64(i2)*100)
		fmt.Println("%\n")
		ColorPrint("[info] 答题结束，对方用时"+strconv.FormatInt(int64(opponent.consume), 10)+"秒，共回答正确：" + strconv.FormatInt(int64(opponent.correct), 10) + "题，正确率",FontColor.light_gray)
		fmt.Printf("%.2f",float64(opponent.correct)/float64(i2)*100)
		fmt.Println("%\n")
		//win,_ :=jsonparser.GetBoolean(finish,"data","badge")
		integral ,_ :=jsonparser.GetInt(finish,"data","integral")

		if integral > 0 {
			ColorPrint("[info] 我方获胜，获得积分" + strconv.FormatInt(int64(integral), 10) + "个\n",FontColor.light_green)
		} else{
			ColorPrint("[info] 对方获胜，获得积分" + strconv.FormatInt(int64(integral), 10) + "个\n",FontColor.light_red)
		}
	}
	fmt.Println("按任意键退出...")
	var input string
	fmt.Scanln(&input)
}
func parseAns(content []byte) string {
	var count int = 0;
	var _Options Options
	jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int,err error) {
		id,_ := jsonparser.GetString(value,"id")
		tit,_ := jsonparser.GetString(value,"title")
		_Options.Option=append(_Options.Option,Option{id,getVisable(tit)})
		count++
	}, "data","options")
	_Options.Count=count
	bytes, _ := json.Marshal(_Options)
	return(string(bytes))
}
func findAns(content []byte,ans map[string]string,act_id string,mode_id string,question_id string) Answer{
	var Ans Answer
	Ans.Activity_id=act_id
	Ans.Way="1"
	Ans.Mode_id=mode_id
	Ans.Question_id=question_id
	jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int,err error) {
		id,_ := jsonparser.GetString(value,"id")
		tit,_ := jsonparser.GetString(value,"title")
		for _,val := range ans {
			if val == getVisable(tit) {
				Ans.Answer=append(Ans.Answer,id)
			}
		}
	}, "data","options")
	return Ans
}
func getVisable(html string)string{
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("[style*='display:none;']").Each(func(i int, selection *goquery.Selection) {
		selection.SetText("")
	})
	doc.Find("[style*='display: none;']").Each(func(i int, selection *goquery.Selection) {
		selection.SetText("")
	})
	return(compressStr(doc.Text()))
}
func _md5(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		log.Fatal(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}
func getAns(content []byte,qid string) {
	var count int = 0;
	var _Options Options
	db, _ := sql.Open("sqlite3", "./question.db")
	var options string
	db.QueryRow("select options from question where `id` = \""+qid+"\"",1).Scan(&options)
	//fmt.Println(options)
	var opmap map[string]string
	opmap=make(map[string]string)
	jsonparser.ArrayEach([]byte(options), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id,err := jsonparser.GetString(value,"Id")
		tit,_ := jsonparser.GetString(value,"Title")
		opmap[id]=tit
	},"Option")
	jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id := string(value)
		tit := opmap[id]
		_Options.Option=append(_Options.Option,Option{id,tit})
		count++
	}, "data","correct_ids")
	_Options.Count=count
	bytes, _ := json.Marshal(_Options)
	//fmt.Println(opmap)
	smst,_:=db.Prepare("UPDATE `question` SET `answer` = ? where `id` = ?")
	//fmt.Println(err)
	smst.Exec(string(bytes),qid)
}
func compressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}