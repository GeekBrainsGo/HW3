/*
Go. Homework 3
Zaur Malakhov, dated Sep 21, 2019
*/

package main

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	r := chi.NewRouter()
	lg := logrus.New()

	serv := Server{
		lg:    lg,
		Title: "БЛОГ О СПОРТЕ",
		Posts: Posts{
			{
				Id:    1,
				Title: "Как быстро набрать мышечную массу",
				Date:  "2019-02-17",
				SmallDescription: "Многие хотят быть большими и сильными, но как это сделать знают не все. " +
					"Спорт – это здоровье, долголетие, привлекательность и активность… Кто из нас не хочет " +
					"приумножить здоровье и привлекательность, жить в полную силу и показать на что способен. " +
					"Выделиться, стать больше очень заманчивая идея, но эволюцией устроено так, ...",
				Description: "Прежде, чем задаваться вопросом как быстро набрать мышечную массу, следует помнить, " +
					"что мышцы растут во время отдыха – микроразрывы мышц восстанавливаются с запасом – это и " +
					"есть прирост массы. Значит, нужно обеспечить отдых, и 8 часовой сон это отличная помощь в " +
					"этом процессе. Перед сном не желательно ужинать за несколько часов, чтобы происходил процесс " +
					"восстановления мышечных волокон, а не пищеварения. В среднем из расчета на килограмм веса нужно " +
					"потреблять 2 грамма белка. Для равномерного отдыха обычно тренируются 3 раза в неделю через день. " +
					"Среди упражнений должна быть база, так называемая золотая тройка – жим лежа, присед, становая " +
					"тяга. Они включают максимально все группы мышц, щедро снабжая тело кровообращением. " +
					"Так как необходимо увеличить мышечную массу, то упражнения должны выполняться в таком " +
					"режиме – больше масса, меньше повторов, больше отдых, в противном случае вы разовьете " +
					"выносливость и силу, без прироста объемов. Если это не помогает, нужно выяснить свою " +
					"конституцию, мышечный тип.",
			},
			{
				Id:    2,
				Title: "Комплекс упражнений для сжигания подкожного жира на животе",
				Date:  "2019-07-15",
				SmallDescription: "Привет всем неравнодушным сторонникам здорового образа жизни и тем, кто верит " +
					"в свои силы для построения тела мечты! В этой статье я уделю внимание главной проблеме " +
					"худеющих, а именно подкожному жиру на животе и способам борьбы с ним. Упражнения для " +
					"сжигания жира на животе при регулярности занятий дадут ожидаемый эффект и наполнят " +
					"энергией для дальнейшего изучения возможностей своего организма.",
				Description: "Действующим способом борьбы с большим животом будет уменьшение количества " +
					"потребляемых углеводов для того чтобы в «топливо» поступал лишний жир и увеличилась " +
					"физическая активность за счет силовых тренировок. База требует много энергии на " +
					"восстановление и окисление мышц, которая берется из жировых запасов. После тяжелой " +
					"тренировки ускоряется метаболизм и запущенный процесс жиросжигания продолжается даже " +
					"во время сна. У девушек из-за генетической и природной предрасположенности к полноте, " +
					"накапливается рыхлый подкожный жир в области ушек, бедер и живота. Для женщин будут " +
					"эффективны прогулки, пробежка, ходьба по лестницам вместо лифта, силовые упражнения " +
					"с небольшим весом и большим количеством повторений, кардио нагрузка (велотренажер, бег, " +
					"скакалка и высокоинтервальные супер-сеты), бодифлекс и фитнес. Для того чтобы разогнать " +
					"лимфу и кровоток в жировых клетках можно подключить массажи с жесткой щеткой перед и " +
					"после тренировки по 5 минут до появления чувства тепла.",
			},
			{
				Id:               3,
				Title:            "Комплекс упражнений для утренней зарядки",
				Date:             "2019-09-12",
				SmallDescription: "В сегодняшней статье предлагаю рассмотреть комплекс упражнений для утренней зарядки...",
				Description: "1. Разминаем шею. Проснувшись утром, не стоит долго лежать в кровати. " +
					"Незамедлительно вставайте, включайте ритмичную музыку и начинайте разогревать мышцы. " +
					"Начало утренней зарядки не должно сопровождаться резкими движениями, иначе увеличиваете " +
					"риск получения неприятного растяжения. Разминка начинается с плавных размеренных наклонов " +
					"головы из стороны в сторону. Повторяйте не менее 10 раз. " +
					"2. Прорабатываем плечи. Теперь задействуем мышцы плечевого сустава. Станьте ровно, " +
					"подняв руки вверх. После этого, не сгибая локти, выполняйте круговые вращения руками, " +
					"напоминающие «ветряную мельницу». Повторяйте около минуты эти активные движения. " +
					"3. Разминка рук. Ноги остаются в том же положении – на ширине плеч. Взяв руки «в замок», " +
					"начинайте медленно вращать кистями по часовой стрелке. Продолжайте выполнять упражнения 30 секунд.",
			},
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleGetIndexHtml)
		r.Get("/view/{postID}", serv.HandleGetPostHtml)
		r.Get("/edit/{postID}", serv.HandleGetEditHtml)
	})

	logrus.Info("server is starts")
	http.ListenAndServe(":8080", r)
}

type Server struct {
	lg    *logrus.Logger
	Title string
	Posts Posts
}

type Posts []Post
type Post struct {
	Id               int
	Title            string
	Date             string
	SmallDescription string
	Description      string
}

func (serv *Server) HandleGetIndexHtml(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", serv)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleGetPostHtml(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)
	postID--

	post := serv.Posts[postID]

	file, _ := os.Open("./www/static/post.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}

func (serv *Server) HandleGetEditHtml(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, _ := strconv.ParseInt(postIDStr, 10, 64)
	postID--

	post := serv.Posts[postID]

	file, _ := os.Open("./www/static/edit.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		serv.lg.WithError(err).Error("template")
	}
}
