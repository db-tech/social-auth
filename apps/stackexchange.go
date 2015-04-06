

package apps

import (
    "github.com/beego/social-auth"
    "github.com/astaxie/beego/httplib"
    "encoding/json"
    "fmt"
    "github.com/astaxie/beego"
    "net/url"
)


type Stackexchange struct {
    BaseProvider
    Key string
}

func (p *Stackexchange) GetType() social.SocialType {
    return social.SocialStackexchange
}

func (p *Stackexchange) GetName() string {
    return "Stackexchange"
}

func (p *Stackexchange) GetPath() string {
    return "stackexchange"
}

func (p *Stackexchange) GetIndentify(tok *social.Token) (string, error) {
    vals := make(map[string]interface{})
    uri := "https://api.stackexchange.com/2.2/me?key="+url.QueryEscape(p.Key)+"&site=stackoverflow&order=asc&sort=name&access_token="+url.QueryEscape(tok.AccessToken)+"&filter=default"
    req := httplib.Get(uri)
    req.SetTransport(social.DefaultTransport)

    resp, err := req.Response()
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    decoder.UseNumber()

    if err := decoder.Decode(&vals); err != nil {
        beego.Error("Get Identify: DecodeError:"+err.Error())
        return "", err
    }
    beego.Debug(vals)
    if vals["error_id"] != nil {
        beego.Error("Get Identify: error_message: ",vals["error_message"])
        return "", fmt.Errorf("%v", vals["error_message"])
    }

    if vals["items"] == nil {
        beego.Error("Get Identify: items==nil")
        return "", nil
    }

    itemsArray := vals["items"].([]interface{})
    itemMap := itemsArray[0].(map[string]interface{})

    if itemMap["account_id"] == nil {
        beego.Error("Get Identify: account_id==nil")
        return "", nil
    }

    return fmt.Sprint(itemMap["account_id"]), nil
}

var _ social.Provider = new(Stackexchange)

func NewStackexchange(clientid, secret string, key string) *Stackexchange {
    p := new(Stackexchange)
    p.App = p
    p.ClientId = clientid
    p.ClientSecret = secret
    p.Scope = ""
    p.AuthURL = "https://stackexchange.com/oauth"
    p.TokenURL = "https://stackexchange.com/oauth/access_token"
    p.RedirectURL = social.DefaultAppUrl + "login/stackexchange/access"
    p.AccessType = "offline"
    p.ApprovalPrompt = "auto"
    p.Key = key
    return p
}