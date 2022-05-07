package model

import (
	"github.com/0xJacky/Nginx-UI/server/tool/nginx"
	"io/ioutil"
	"path/filepath"
)

type Cert struct {
	Model
	Domain string `json:"domain"`
}

func FirstCert(domain string) (c Cert, err error) {
	err = db.First(&c, &Cert{
		Domain: domain,
	}).Error

	return
}

func FirstOrCreateCert(domain string) (c Cert, err error) {
	err = db.FirstOrCreate(&c, &Cert{Domain: domain}).Error
	return
}

func GetAutoCertList() (c []Cert) {
	var t []Cert
	db.Find(&t)
	// check if this domain is enabled

	enabledConfig, err := ioutil.ReadDir(filepath.Join(nginx.GetNginxConfPath("sites-enabled")))

	if err != nil {
		return
	}

	enabledConfigMap := make(map[string]bool)
	for i := range enabledConfig {
		enabledConfigMap[enabledConfig[i].Name()] = true
	}

	for _, v := range t {
		if enabledConfigMap[v.Domain] == true {
			c = append(c, v)
		}
	}

	return
}

func (c *Cert) Remove() error {
	return db.Where("domain", c.Domain).Delete(c).Error
}
