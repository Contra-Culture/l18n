package l18n_test

import (
	"errors"

	. "github.com/Contra-Culture/l18n"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("l18n", func() {
	Describe("L18n", func() {
		Describe("L()", func() {
			It("returns *L18n", func() {
				l := L([]string{"en", "ru", "ua"})
				Expect(l).NotTo(BeNil())
			})
		})
		Describe(".Add()", func() {
			Context("when all translations provided", func() {
				It("adds translations", func() {
					l := L([]string{"en", "ru", "ua"})
					err := l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
						"en": "home",
						"ru": "главная",
						"ua": "головна",
					})
					Expect(err).NotTo(HaveOccurred())
					err = l.Add([]string{"main", "nav", "archive"}, map[string]interface{}{
						"en": "archive",
						"ru": "архив",
						"ua": "ахів",
					})
					Expect(err).NotTo(HaveOccurred())
				})
			})
			Context("when not all translations provided", func() {
				It("fails and returns error", func() {
					l := L([]string{"en", "ru", "ua"})
					err := l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
						"en": "home",
						"ru": "главная",
					})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("translation for \"ua\" language is not provided"))
				})
			})
			Context("when there is a conflict between translation and scope", func() {
				It("fails and returns error", func() {
					l := L([]string{"en", "ru", "ua"})
					l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
						"en": "home",
						"ru": "главная",
						"ua": "головна",
					})
					err := l.Add([]string{"main", "nav"}, map[string]interface{}{
						"en": "navigation",
						"ru": "навигация",
						"ua": "навігація",
					})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("wrong path \"main/nav\": there is already exists a scope"))
				})
			})
			Context("when translation already exists", func() {
				It("fails and returns error", func() {
					l := L([]string{"en", "ru", "ua"})
					l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
						"en": "home",
						"ru": "главная",
						"ua": "головна",
					})
					err := l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
						"en": "home",
						"ru": "главная",
						"ua": "головна",
					})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("translation \"main/nav/home\" already exists"))
				})
			})
		})
		Describe(".Lang()", func() {
			Context("when language exists", func() {
				It("returs language-scoped context", func() {
					l := L([]string{"en", "ru", "ua"})
					l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
						"en": "home",
						"ru": "главная",
						"ua": "головна",
					})
					scoped, err := l.Lang("en")
					Expect(err).NotTo(HaveOccurred())
					Expect(scoped).NotTo(BeNil())
				})
			})
			Context("when language doesn't exist", func() {
				It("fails and returns error", func() {
					l := L([]string{"en", "ru", "ua"})
					scoped, err := l.Lang("fr")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("language \"fr\" is not registered"))
					Expect(scoped).To(BeNil())
				})
			})
		})
		Describe("Scoped", func() {
			Describe(".Get()", func() {
				Context("when translation exists", func() {
					It("returns translations", func() {
						l := L([]string{"en", "ru", "ua"})
						l.Add([]string{"main", "nav", "home"}, map[string]interface{}{
							"en": "home",
							"ru": "главная",
							"ua": "головна",
						})
						scoped, _ := l.Lang("en")
						v, err := scoped.Get([]string{"main", "nav", "home"}, map[string]interface{}{})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("home"))
						l.Add([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"en": func(args map[string]interface{}) (string, error) {
								rawSex, ok := args["sex"]
								if !ok {
									return "", errors.New("\"sex\" is not provided")
								}
								sex, ok := rawSex.(string)
								if !ok {
									return "", errors.New("\"sex\" should be a string")
								}
								switch sex {
								case "woman":
									return "she applied your invitation", nil
								case "man":
									return "he applied your invitation", nil
								default:
									return "your invitation was applied", nil
								}
							},
							"ru": func(args map[string]interface{}) (string, error) {
								rawSex, ok := args["sex"]
								if !ok {
									return "", errors.New("\"sex\" is not provided")
								}
								sex, ok := rawSex.(string)
								if !ok {
									return "", errors.New("\"sex\" should be a string")
								}
								switch sex {
								case "woman":
									return "приняла ваше приглашение", nil
								case "man":
									return "принял ваше приглашение", nil
								default:
									return "ваше приглашение было принято", nil
								}
							},
							"ua": func(args map[string]interface{}) (string, error) {
								rawSex, ok := args["sex"]
								if !ok {
									return "", errors.New("\"sex\" is not provided")
								}
								sex, ok := rawSex.(string)
								if !ok {
									return "", errors.New("\"sex\" should be a string")
								}
								switch sex {
								case "woman":
									return "прийняла ваше запрошення", nil
								case "man":
									return "прийняв ваше запрошення", nil
								default:
									return "ваше запрошення було прийняте", nil
								}
							},
						})
						scoped, _ = l.Lang("en")
						v, err = scoped.Get([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"sex": "man",
						})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("he applied your invitation"))
						v, err = scoped.Get([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"sex": "woman",
						})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("she applied your invitation"))
						scoped, _ = l.Lang("ru")
						v, err = scoped.Get([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"sex": "man",
						})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("принял ваше приглашение"))
						v, err = scoped.Get([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"sex": "woman",
						})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("приняла ваше приглашение"))
						scoped, _ = l.Lang("ua")
						v, err = scoped.Get([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"sex": "man",
						})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("прийняв ваше запрошення"))
						v, err = scoped.Get([]string{"main", "invitation", "applied"}, map[string]interface{}{
							"sex": "woman",
						})
						Expect(err).NotTo(HaveOccurred())
						Expect(v).To(Equal("прийняла ваше запрошення"))
					})
				})
				Context("when translation doesn't exist", func() {
					It("fails and returns error", func() {
						l := L([]string{"en", "ru", "ua"})
						l.Add([]string{"main", "nav", "home"}, map[string]interface{}{})
						scoped, _ := l.Lang("en")
						v, err := scoped.Get([]string{"main", "nav", "home"}, map[string]interface{}{})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("translation \"main/nav/home\" does not exist"))
						Expect(v).To(Equal(""))
					})
				})
			})
		})
	})
})
