package yuwiki

import (
	"errors"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"time"
)

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpForm struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type CaptchaResult struct {
	CaptchaId  string `json:"captchaId"`
	CaptchaUrl string `json:"captchaUrl"`
}

type UserProfile struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"required"`
	Gender   int8   `json:"gender"`
	Birthday string `json:"birthday"`
}

type PageDto struct {
	ID        uint     `json:"id"`
	BookId    uint     `json:"bookId" binding:"required"`
	PartId    uint     `json:"partId" binding:"required"`
	Title     string   `json:"title" binding:"required"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Owner     uint     `json:"owner"`
	SortCode  uint     `json:"sortCode"`
	Published bool
	CreatedAt time.Time `json:"CreatedAt"`
}

type UserSearch struct {
	ID      uint   `json:"id"`
	Keyword string `json:"keyword"`
	Gender  int8   `json:"gender"`
	NoSelf  bool   `json:"noSelf"`
}

type PasswordModify struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type SortBook struct {
	BookId   uint `json:"bookId"`
	SortCode uint `json:"sortCode"`
}

type SortPart struct {
	PartId   uint `json:"partId"`
	SortCode uint `json:"sortCode"`
}

type SortPage struct {
	PageId   uint `json:"pageId"`
	SortCode uint `json:"sortCode"`
}

func checkLogin(c *gin.Context) (bool, *User, error) {
	loginForm := &LoginForm{}
	if err := c.ShouldBind(loginForm); err != nil {
		return false, nil, err
	}
	if len(loginForm.Password) < 6 {
		return false, nil, errors.New("密码长度不能小于6位")
	}
	user := &User{}
	if err := Db.Where("username = ?", loginForm.Username).Find(user).Error; err != nil {
		return false, nil, errors.New(fmt.Sprintf("用户 %s 不存在", loginForm.Username))
	}
	if Match(loginForm.Password, user.Salt, user.Password) {
		return true, user, nil
	}
	return false, nil, nil
}

func signUpHandler(c *gin.Context) {
	signUpForm := &SignUpForm{}
	if err := c.ShouldBind(signUpForm); err != nil {
		Result(c, CodeFail(ParamsError))
		return
	}
	emailReg := regexp.MustCompile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`)
	if !emailReg.MatchString(signUpForm.Email) {
		Result(c, MsgFail("邮箱格式有误"))
		return
	}
	if signUpForm.Password != signUpForm.ConfirmPassword {
		Result(c, MsgFail("两次输入的密码不一致"))
		return
	}
	dbUser := &User{}
	Db.Where("username = ?", signUpForm.Username).Find(dbUser)
	if dbUser.ID != 0 {
		Result(c, MsgFail(fmt.Sprintf("用户 %s 已存在", signUpForm.Username)))
		return
	}
	salt := GenSalt()
	user := &User{
		Username:        signUpForm.Username,
		Email:           signUpForm.Email,
		InitPassword:    "",
		PasswordChanged: true,
		Salt:            salt,
	}
	user.Password, _ = EncPassword(signUpForm.Password, salt)
	if err := Db.Create(user).Error; err != nil {
		Result(c, MsgFail("注册账号失败"))
	} else {
		Result(c, Ok())
	}
}

func getCaptchaHandler(c *gin.Context) {
	if captchaId := captcha.New(); captchaId != "" {
		captchaResult := &CaptchaResult{
			CaptchaId:  captchaId,
			CaptchaUrl: "/captcha/show/" + captchaId + ".png",
		}
		Result(c, OkData(captchaResult))
	} else {
		Result(c, CodeFail(FAIL))
	}
}

func showCaptchaHandler(c *gin.Context) {
	source := c.Param("source")
	log.Debugf("获取验证码图片：%s", source)
	ServeHTTP(c.Writer, c.Request)
}

func verifyCaptchaHandler(c *gin.Context) {
	captchaId := c.Query("captchaId")
	value := c.Query("value")
	if captchaId == "" || value == "" || len(value) != 6 {
		Result(c, CodeFail(ParamsError))
	} else if captcha.VerifyString(captchaId, value) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(CaptchaVerifyFail))
	}
}

func getBooksHandler(c *gin.Context) {
	Result(c, OkData(getBooks()))
}

func getBookPartsHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getBookParts(uint(bookId))))
	}
}

func getPartHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPart(uint(partId))))
	}
}

func saveBookHandler(c *gin.Context) {
	book := &Book{}
	if err := c.ShouldBindJSON(book); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if saveBook(book) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func deleteBookHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deleteBook(uint(bookId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func shareBookHandler(c *gin.Context) {
	sharedBook := &SharedBook{}
	if err := c.ShouldBind(sharedBook); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if saveSharedBook(sharedBook) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func cancelShareBookHandler(c *gin.Context) {
	sharedBook := &SharedBook{}
	if err := c.ShouldBind(sharedBook); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deleteSharedBook(sharedBook) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func getBookSharedUsersHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getBookSharedUsers(uint(bookId))))
	}
}

func getPartPagesHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPartPages(uint(partId))))
	}
}

func savePartHandler(c *gin.Context) {
	part := &Part{}
	if err := c.ShouldBind(part); err != nil {
		Result(c, CodeFail(ParamsError))
		return
	}
	if part.Protected && part.Password == "" {
		Result(c, CodeFail(ParamsError))
		return
	}
	if ok, err := savePart(part); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func deletePartHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deletePart(uint(partId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func getPageHandler(c *gin.Context) {
	pageId, pageIdErr := strconv.ParseUint(c.Param("pageId"), 10, 32)
	editable, editableErr := strconv.ParseBool(c.DefaultQuery("editable", "false"))
	if pageIdErr != nil || editableErr != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPage(uint(pageId), editable)))
	}
}

func savePageHandler(c *gin.Context) {
	pageDto := &PageDto{}
	if err := c.ShouldBind(pageDto); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if savePage(pageDto) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func editPageHandler(c *gin.Context) {
	pageDto := &PageDto{}
	if err := c.ShouldBind(pageDto); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if editPage(pageDto) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(UpdateFail))
	}
}

func deletePageHandler(c *gin.Context) {
	if pageId, err := strconv.ParseUint(c.Param("pageId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deletePage(uint(pageId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func getHistoricalPagesHandler(c *gin.Context) {
	if pageId, err := strconv.ParseUint(c.Param("pageId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getHistoricalPages(uint(pageId))))
	}
}

func getTagsHandler(c *gin.Context) {
	Result(c, OkData(getTags()))
}

func getTemplatesHandler(c *gin.Context) {
	Result(c, OkData(getPageTemplates()))
}

func getTemplateHandler(c *gin.Context) {
	if templateId, err := strconv.ParseUint(c.Param("templateId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getPageTemplate(uint(templateId))))
	}
}

func saveTemplateHandler(c *gin.Context) {
	template := &PageTemplate{}
	if err := c.ShouldBindJSON(template); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if savePageTemplate(template) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func deleteTemplateHandler(c *gin.Context) {
	if templateId, err := strconv.ParseUint(c.Param("templateId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if deletePageTemplate(uint(templateId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(DeleteFail))
	}
}

func getUserInfoHandler(c *gin.Context) {
	if user, err := getCurrentUser(); err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, OkData(user))
	}
}

func editUserHandler(c *gin.Context) {
	userProfile := &UserProfile{}
	if err := c.ShouldBind(userProfile); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if user, err := currentUser(); err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		user.Avatar = userProfile.Avatar
		user.Nickname = userProfile.Nickname
		user.Phone = userProfile.Phone
		user.Email = userProfile.Email
		user.Gender = userProfile.Gender
		user.Birthday, _ = time.ParseInLocation(DateTimeLayout, userProfile.Birthday+" 00:00:00", time.Local)
		if err := Db.Save(user).Error; err != nil {
			Result(c, CodeFail(UpdateFail))
		} else {
			Result(c, Ok())
		}
	}
}

func modifyPasswordHandler(c *gin.Context) {
	modify := &PasswordModify{}
	if err := c.ShouldBind(modify); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, msg := modifyPassword(modify); ok {
		Result(c, Ok())
	} else {
		Result(c, MsgFail(msg))
	}
}

func searchUsersHandler(c *gin.Context) {
	userSearch := &UserSearch{}
	if err := c.ShouldBind(userSearch); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(searchUsers(userSearch)))
	}
}

func siteSearchHandler(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	if keyword == "" {
		Result(c, Ok())
	} else {
		Result(c, OkData(searchPagesByKeyword(keyword)))
	}
}

func sortBooksHandler(c *gin.Context) {
	sortedBooks := &[]SortBook{}
	if err := c.ShouldBind(sortedBooks); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, err := sortBooks(sortedBooks); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SortFail))
	}
}

func sortPartsHandler(c *gin.Context) {
	sortedParts := &[]SortPart{}
	if err := c.ShouldBind(sortedParts); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, err := sortParts(sortedParts); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SortFail))
	}
}

func sortPagesHandler(c *gin.Context) {
	sortedPages := &[]SortPage{}
	if err := c.ShouldBind(sortedPages); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if ok, err := sortPages(sortedPages); ok {
		Result(c, Ok())
	} else if err != nil {
		Result(c, MsgFail(err.Error()))
	} else {
		Result(c, CodeFail(SortFail))
	}
}

func toggleStarBookHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if toggleStarBook(uint(bookId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func toggleStarPartHandler(c *gin.Context) {
	if partId, err := strconv.ParseUint(c.Param("partId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if toggleStarPart(uint(partId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func toggleStarPageHandler(c *gin.Context) {
	if pageId, err := strconv.ParseUint(c.Param("pageId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else if toggleStarPage(uint(pageId)) {
		Result(c, Ok())
	} else {
		Result(c, CodeFail(SaveFail))
	}
}

func getSharedBooksHandler(c *gin.Context) {
	Result(c, OkData(getSharedBooks()))
}

func getSharedPartsHandler(c *gin.Context) {
	if bookId, err := strconv.ParseUint(c.Param("bookId"), 10, 32); err != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getSharedParts(uint(bookId))))
	}
}

func getSharedPagesHandler(c *gin.Context) {
	bookId, bookIdErr := strconv.ParseUint(c.Param("bookId"), 10, 32)
	partId, partIdErr := strconv.ParseUint(c.Param("partId"), 10, 32)
	if bookIdErr != nil || partIdErr != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getSharedPages(uint(bookId), uint(partId))))
	}
}

func viewSharedPageHandler(c *gin.Context) {
	bookId, bookIdErr := strconv.ParseUint(c.Param("bookId"), 10, 32)
	partId, partIdErr := strconv.ParseUint(c.Param("partId"), 10, 32)
	pageId, pageIdErr := strconv.ParseUint(c.Param("pageId"), 10, 32)
	if bookIdErr != nil || partIdErr != nil || pageIdErr != nil {
		Result(c, CodeFail(ParamsError))
	} else {
		Result(c, OkData(getSharedPage(uint(bookId), uint(partId), uint(pageId))))
	}
}
