package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

//文章结构体
type Article struct {
	Id       int64
	Title    string    `orm:"size(64);unique"`
	Category string    `orm:"size(32)"`
	Content  string    `orm:"size(5000)"`
	Created  time.Time `orm:"index;auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"index;auto_now;type(datetime)"`
	Views    int64     `orm:"index"`
}

//自定义表名
func (this *Article) TableName() string {
	return "article"
}

//自定义引擎
func (this *Article) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModel(new(Article))
}

//添加文章
func AddArticle(title, category, content string) error {
	o := orm.NewOrm()
	article := &Article{
		Title:    title,
		Category: category,
		Content:  content,
		Updated:  time.Now(),
	}
	created, _, err := o.ReadOrCreate(article, "Name")
	if err != nil {
		return err
	}
	if created {
		return nil
	} else {
		return errors.New("文章标题已经存在")
	}
	//更新分类下面的文章数量
	cate := &Category{Name: category}
	err = o.Read(cate, "Name")
	if err != nil {
		cate.Count++
		err = o.Update(cate, "Count")
	}
	return err
}

//修改文章
func EditArticle(tid, title, category, content string) error {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	article := &Article{Id: id}
	var oldCategory string
	if o.Read(article) == nil {
		if article.Category != category {
			oldCategory = article.Category
		}
		article.Title = title
		article.Category = category
		article.Content = content
		_, err = o.Update(article)
		if err != nil {
			return err
		}
	}
	if len(oldCategory) > 0 {
		cate := new(Category)
		err = o.QueryTable("category").Filter("name", oldCategory).One(cate)
		if err != nil {
			return err
		}
		cate.Count--
		_, err = o.Update(cate)
	}

	cate := new(Category)
	err = o.QueryTable("category").Filter("name", category).One(cate)
	if err == nil {
		cate.Count++
		_, err = o.Update(cate)
	}
	return err
}

//删除文章
func DelArticle(aid string) error {
	id, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		return nil
	}
	o := orm.NewOrm()
	article := &Article{Id: id}
	var oldCategory string
	if o.Read(article) == nil {
		oldCategory = article.Category
		_, err = o.Delete(article)
		if err != nil {
			return err
		}
	}
	if len(oldCategory) > 0 {
		cate := new(Category)
		err = o.QueryTable("category").Filter("name", oldCategory).One(cate)
		if err != nil {
			return err
		}
		cate.Count--
		_, err = o.Update(cate)
	}
	return err
}

//获取一篇文章
func GetArticle(aid string) (*Article, error) {
	id, err := strconv.ParseInt(aid, 10, 64)
	if err != nil {
		return err
	}
	article := &Article{Id: id}
	o := orm.NewOrm()
	err = o.Read(article)
	if err != nil {
		return err
	}
	article.Views++
	_, err = o.Update(article)
	return article, err
}

//获取文章列表
func GetArticles(start, offset, string) ([]*Article, error) {
	articles := make([]*Article, 0)
	o := orm.NewOrm()
}