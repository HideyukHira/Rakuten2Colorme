package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"net/http"
	"strings"
)

const (
	CMA = ","
	NN  = ""
)

type Product struct {
	serial  string //3	型番
	name    string //4	"商品名*（必須）"
	explain string //34	商品説明
	price   string //16	販売価格
	zaiko   string //21	在庫管理「0」：在庫管理する	「1」：在庫管理しない
	show    string //56	"掲載設定*（必須）"
}


func main() {
	//コマンドライン引数 大カテゴリ URL
	flag.Parse()
	args := flag.Args()

	var categoryname = args[0] //1	"カテゴリー（大）*（必須）"
	var categoryListPage = args[1]

	//カテゴリ一覧から商品一覧URL取得
	urllists := []string{}
	urllists = Pagelist(categoryListPage)

	//商品一覧URLをループしてCSV用のデータを標準出力
	for _, urlvalue := range urllists {

		Itemscrape(categoryname, urlvalue)
	}

}

func Pagelist(categoryListPage string) []string {

	//カテゴリ一覧から商品URLのリストを返す
	resp, err := http.Get(categoryListPage)

	if err != nil {
	}
	defer resp.Body.Close()

	//楽天はEUC-JPなのでエンコードする
	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.EUCJP.NewDecoder())
	//ページ読み込み

	doc, err := goquery.NewDocumentFromReader(utfBody)

	//型番
	urllists := []string{}
	doc.Find(".category_itemnamelink").Each(func(i int, s *goquery.Selection) {
		var url, _ = s.Attr("href")
		urllists = append(urllists, url)
	})

	return urllists

}

func Itemscrape(categoryname string, url string) {
	//個別アイテムのスクレイピングと標準出力

	var p Product
	/*
		初期設定

	*/

	//在庫フラグ 在庫管理「0」：在庫管理する	「1」：在庫管理しない
	p.zaiko = "1"
	//掲載フラグ 掲載する
	p.show = "0"

	resp, err := http.Get(url)

	if err != nil {
	}
	defer resp.Body.Close()

	//楽天はEUC-JPなのでエンコードする
	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.EUCJP.NewDecoder())
	//ページ読み込み

	doc, err := goquery.NewDocumentFromReader(utfBody)

	//型番
	doc.Find("span.item_number").Each(func(i int, s *goquery.Selection) {
		p.serial = s.Text()
	})

	//商品名
	doc.Find("span.item_name").Each(func(i int, s *goquery.Selection) {
		p.name = s.Text()
	})

	//　商品画像シリーズ
	var imgs [10]string
	doc.Find(".rakutenLimitedId_ImageMain1-3").Each(func(i int, s *goquery.Selection) {
		var imgurl, _ = s.Attr("href")
		imgs[i] = imgurl
	})

	//商品説明  フィールドを"で囲むこと。"は""でエスケープすること
	doc.Find("span.sale_desc").Each(func(i int, s *goquery.Selection) {
		var explain, _ = s.Html()
		explain = strings.Replace(explain, "\"", "\"\"", -1)
		p.explain = "\"" + explain + "\""
	})

	//販売価格
	doc.Find("input#ratPrice").Each(func(i int, s *goquery.Selection) {
		var str, _ = s.Attr("value")
		p.price = str
	})

	fmt.Print(
		//1	"カテゴリー（大）*（必須）"
		categoryname, CMA,
		//2	カテゴリー（小）
		NN, CMA,
		//3	型番
		p.serial, CMA,
		//4	"商品名*（必須）"
		p.name, CMA,
		//5	商品画像
		imgs[0], CMA,
		//6	モバイルショップ用商品画像作成
		NN, CMA,
		//7	その他画像1
		imgs[1], CMA,
		//8	その他画像2
		imgs[2], CMA,
		//9	その他画像3
		imgs[3], CMA,
		//10	その他画像4
		imgs[4], CMA,
		//11	その他画像5
		imgs[5], CMA,
		//12	その他画像6
		imgs[6], CMA,
		//13	その他画像7
		imgs[7], CMA,
		//14	その他画像8
		imgs[8], CMA,
		//15	その他画像9
		imgs[9], CMA,
		//16	販売価格
		p.price, CMA,
		//17	会員価格
		NN, CMA,
		//18	定価
		NN, CMA,
		//19	原価
		NN, CMA,
		//20	在庫数
		NN, CMA,
		//21	在庫管理
		p.zaiko, CMA,
		//22	最小購入数量
		NN, CMA,
		//23	最大購入数量
		NN, CMA,
		//24	販売開始日付
		NN, CMA,
		//25	販売開始時間
		NN, CMA,
		//26	販売終了日付
		NN, CMA,
		//27	販売終了時間
		NN, CMA,
		//28	単位
		NN, CMA,
		//29	重量
		NN, CMA,
		//30	売切れ時の 表示設定
		NN, CMA,
		//31	適正在庫数
		NN, CMA,
		//32	表示順
		NN, CMA,
		//33	簡易説明
		NN, CMA,
		///34	商品説明
		p.explain, CMA,
		//35	モバイルショップ用商品説明
		NN, CMA,
		//36	スマートフォンショップ用商品説明
		NN, CMA,
		//37	未使用項目
		NN, CMA,
		//38	Newマーク付加設定
		NN, CMA,
		//39	Newマーク画像選択
		NN, CMA,
		//40	カラメル用カテゴリーID
		NN, CMA,
		//41	カラメル用タグ1
		NN, CMA,
		//42	カラメル用タグ2
		NN, CMA,
		//43	カラメル用タグ3
		NN, CMA,
		//44	カラメル用商品説明
		NN, CMA,
		//45	ブランド [固有商品ID]
		NN, CMA,
		//46	JAN/ISBN (GTIN) [固有商品ID]
		NN, CMA,
		//47	製品番号 (MPN) [固有商品ID]
		NN, CMA,
		//48	状態
		NN, CMA,
		//49	性別
		NN, CMA,
		//50	色
		NN, CMA,
		//51	サイズ
		NN, CMA,
		//52	タイトル要素
		NN, CMA,
		//53	meta要素 キーワード設定
		NN, CMA,
		//54	meta要素 ページ概要
		NN, CMA,
		//55	個別送料
		NN, CMA,
		//56	"掲載設定*（必須）"
		p.show+"\n",
	)
}
