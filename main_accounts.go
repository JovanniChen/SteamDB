package main

var accounts = []Account{
	{Username: "zszvlv6362", Password: "ejuj7Rnof1BB", SharedSecret: "mQI147JxRz78GWjDdQEBoL7aaBc="},   // [0] [45]
	{Username: "gbmqnl7210", Password: "i80sMCigz1rw", SharedSecret: "Uinb4sxNpcP8KQBcYgdAZ2eiJDg="},   // [1] [46]
	{Username: "zdckla1506", Password: "d3c9InY7Epwi", SharedSecret: "rQJ4b42FyGsvGcp6XYx+SEYylyo="},   // [2] [47]
	{Username: "uvtjrm4501", Password: "u9NIlsVugLH5", SharedSecret: "y77Jk5v4rxrck/149zDMB+b3s/U="},   // [3] [48]
	{Username: "ddndd12412", Password: "New0KJYVv16", SharedSecret: "VoSY5VrnD+CJooEVrlADofTGTok="},    // [4] [51]
	{Username: "ttmsq72777", Password: "yoRD7x6LQvgu", SharedSecret: "5boHTiGFhQoszGcpFDLB7H7thng="},   // [5] [52]
	{Username: "xqkea03549", Password: "wuwQJ5WFdZp1", SharedSecret: "59z0KMWJFdgfWrSgYYADD/LBPyU="},   // [6] [53]
	{Username: "ffotd74229", Password: "oP4M4CMHAftX", SharedSecret: "IDhBX3NM+8fZCti4C3d6oFhXI6E="},   // [7] [54]
	{Username: "j47hoord6j", Password: "NewRP7IhC9Z", SharedSecret: "Gwgztog4anK0soQp4IgLaZIki0s="},    // [8] [57] 市场不可用
	{Username: "naotqp7801", Password: "ja9C5LZelku0", SharedSecret: "g+kIH7JuL98R5O00j87379CkFus="},   // [9]
	{Username: "zszvlv6362", Password: "ejuj7Rnof1BB", SharedSecret: "mQI147JxRz78GWjDdQEBoL7aaBc="},   // [10]
	{Username: "fbrdz08225", Password: "NewNWnME1R6", SharedSecret: "VjYAPygKL4jxwSu69HeyzW58r3M="},    // [11]
	{Username: "rwfio67235", Password: "JzBvNCICYfFx", SharedSecret: "0C4hU7ieyVyYFvdDPKoTII20xMc="},   // [12]
	{Username: "ejvp732231", Password: "myz2bzwCzFYQ", SharedSecret: "KHzBIonDKW8enmoCUYgLN+oYQ4M="},   // [13]
	{Username: "mcg9ipxd04nl", Password: "AAlLQXPdDy3U", SharedSecret: "zUN+RyvAQZjHnyT+5guHBPB2NOg="}, // [14]
	{Username: "pfze6stttee", Password: "P4A9ydYvzGmq", SharedSecret: "qNSLkP8OjsD/VuHG5eFGUjSupCs="},  // [15]
	{Username: "krqk10ik5qk", Password: "vPfPgdUqVX76", SharedSecret: "1daUVnJtxNg4pB2gxV2l10jfz1U="},  // [16]
	{Username: "bhrcrnulng5", Password: "YF8TX9fAxWsq", SharedSecret: "X1z0h/KTJmns1um4ThZdRGrrNps="},  // [17]
	{Username: "dih9u8nad", Password: "BJdENppgNHxH", SharedSecret: "TFuheV7W4oPoH4Q2EH8EEi9vmKU="},    // [18]
	{Username: "tvyij7pxdasz", Password: "0DLuZvp5MSEI", SharedSecret: "YGnTkbpo/uOGFtNeFRhGlIQxrEg="}, // [19]
	{Username: "hnysh898sg", Password: "mC43y8o8irxT", SharedSecret: "ynrDLQop7KGLFe0DfiFcW8lOy6A="},   // [20]
	{Username: "tzjn5e5xnz5y", Password: "wBvFtVZkCpsZ", SharedSecret: "RGcG25ZSJiswJpwz56DTaPWR+nI="}, // [21]
	{Username: "tgf21ra7e3", Password: "50Z2DJJMNnfI", SharedSecret: "FWhvNXuMuhPGj4V1FLBD31Fzks8="},   // [22]
	{Username: "acjjzz1twx", Password: "H3s5wvDLgL4d", SharedSecret: "+xtZKqLMsIMu7T4LgP3rO6wNV2Q="},   // [23]
	{Username: "akji2nh66u0f", Password: "t88j8RZs37BY", SharedSecret: "KjTGnop2fyCu2E7hBWflE2TdEO4="}, // [24]
	{Username: "kneao1dmw", Password: "5Dtu8Kk9JIAg", SharedSecret: "GjkmsGH1utG1QVyg06NQc1wNTwg="},    // [25]
	{Username: "plzhhgt075", Password: "ESvgOZnLTjKb", SharedSecret: "WxqsqlawbHEB7Pjah0wutOpLKE0="},   // [26]
	{Username: "iuuwhmusxdv7", Password: "STQI7NOal7l6", SharedSecret: "+xtZKqLMsIMu7T4LgP3rO6wNV2Q="}, // [27]
	{Username: "bmlgbjot5hz", Password: "1Z5pkOTuZigf", SharedSecret: "GZKHLVfxwjYPFMBF33l7Vu3NMY4="},  // [28]
	{Username: "aeuybz0905", Password: "f4J5Cs6cHnHP", SharedSecret: "UGDYQfigAc47yH/wPcL0E3PCHPY="},   // [29]
	{Username: "maantzmeesnw", Password: "NewQa05drZf", SharedSecret: "8dYeP5BMZ3L70vvwxSubGpHD2bo="},  // [30]
	{Username: "yfsqqkd80", Password: "2J76AxYA6Pbi", SharedSecret: "xA8vLi8mw7kOvghspCYy+J124Fk="},    // [31]
	{Username: "qdkfunqkvffp", Password: "e6UnHb9ydRGv", SharedSecret: "fv5mQVGjQ0GwrgTjwm2pGts/6RA="}, // [32]
	{Username: "fqlyt822", Password: "sfmqjRqOpf4L", SharedSecret: "PdRJRMPqe+D02gd2YmTFtagQMkQ="},     // [33]
	{Username: "dqx1pil5mpj", Password: "F6eyXrLU0XP3", SharedSecret: "AwqrvNq2qxAyR4sNB9Gef9dIGok="},  // [34]
	{Username: "gsmut047", Password: "1oq0aF7mL0pl", SharedSecret: "MWWC3d47Cy7yGxvJQmH+py1KI+I="},     // [35]
	{Username: "gxruaknmby", Password: "SK8rvcMSLq8m", SharedSecret: "sbbo5c0NDKLyDvT6DlF4EbhG3rI="},   // [36]
	{Username: "firby644", Password: "UJZS5JfysEyA", SharedSecret: "I2yFq42TsWb216JfJZUT4IbIbqg="},     // [37]
	{Username: "za0ww9ml4xl2", Password: "HLHxGyRMm6Zi", SharedSecret: "F54xOr9Tpyd5fAxgKx+RHR7vHik="}, // [38]
	{Username: "kikiixbfo", Password: "kMyCmQPiVG", SharedSecret: "MeY/00fytzx5CS8qo1wQ/tHFWv0="},      // [39]
	{Username: "fv5928894", Password: "zim590238", SharedSecret: "4sYqKlQWD4MiZZKUlaFz3JhHhqM="},       // [40]
	{Username: "rwc3665651", Password: "qEOmc2Y6OjuI", SharedSecret: "bW4MbbVQca1ypuLv7IbI2KNAVu8="},   // [41]
}

type Account struct {
	Username     string // Steam用户名
	Password     string // Steam密码
	SharedSecret string // Steam Guard共享密钥(base64编码)
}

func (a *Account) GetUsername() string {
	return a.Username
}
func (a *Account) GetPassword() string {
	return a.Password
}
func (a *Account) GetSharedSecret() string {
	return a.SharedSecret
}

func getAccount(index int) *Account {
	return &accounts[index]
}
