/* Copyright © 2021 Aadhyarupam Innovators  <https://digital.aadhyarupam.com>
 * Source code of "go-ail-examples" is provided under GNU General Public License v3.0
 *
 * Purpose of this project is to demonstrate the code examples of Go (Golang) program,
 * to develop simple, scalable and high performant microservices.

 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
// AppConfig load configuration from file and remote provider (consul or firestore).
// It give preference to remote provider while reading the configuration value if exist.

package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	CONFIG_FILENAME = "appconfig"
	KEY_PROJECT_ID  = "projectid"

	// Remote keystore provider properties
	KEY_REMOTE_PROVIDER = "remote.provider"
	KEY_REMOTE_ENDPOINT = "remote.endpoint"
	KEY_REMOTE_PATH     = "remote.path"

	// SMTP Server Properties
	KEY_SMTP_HOST     = "smtp.host"
	KEY_SMTP_PORT     = "smtp.port"
	KEY_SMTP_IDENTITY = "smtp.identity"
	KEY_SMTP_USER     = "smtp.user"
	KEY_SMTP_SECRET   = "smtp.secret"
	KEY_SMTP_SENDER   = "smtp.sender"

	KEY_OAUTH_GOOGLE_CLIENT_ID     = "oauth.google.clientid"
	KEY_OAUTH_GOOGLE_CLIENT_SECRET = "oauth.google.secret"
	KEY_OAUTH_GOOGLE_REDIRECT_URI  = "oauth.google.redirecturi"

	KEY_SUBSCRIPTION_SERVICE_URL = "subscription-url"

	KEY_HTTP_SERVICE_PORT             = "service.http.port"
	KEY_HTTP_SERVICE_CONTENT_LOCATION = "service.http.contentlocation"

	KEY_DISCOVERY_SERVICE_LOCATION = "service.discovery.location"

	KEY_LOG_LEVEL_DEBUG = "log.debug"
	KEY_LOG_LEVEL_INFO  = "log.info"
	KEY_LOG_LEVEL_TRACE = "log.trace"
)

var configstore *viper.Viper
var useConfigStore bool = false

func AppInit() {
	log.Printf("Loading configuration ..")
	viper.SetEnvPrefix("ail")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("json")
	viper.SetConfigName(CONFIG_FILENAME)
	viper.AddConfigPath("$HOME/.config") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.AddConfigPath("resources")     // path to look for the config file in

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error occured while loading configuration: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found;
			log.Fatalf("Configuration file %s not found, Stopping application startup", CONFIG_FILENAME)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Failed to load configuration file %s, Stopping application startup", CONFIG_FILENAME)
		}
	}
	if viper.IsSet("application") {
		log.Printf("Successfully loaded configuration file: %s", CONFIG_FILENAME)
	}

	// Add Remote provider if exist
	remoteProvider := viper.GetString(KEY_REMOTE_PROVIDER)
	remoteEndpoint := viper.GetString(KEY_REMOTE_ENDPOINT)
	remotePath := viper.GetString(KEY_REMOTE_PATH)
	configstore = viper.New()
	if remoteProvider != "" && remoteEndpoint != "" && remotePath != "" {
		// Don't set ConfigType and ConfigName while using firestore remote provider otherwise you might encounter unexpected end of JSON error
		//configstore.SetConfigName(CONFIG_FILENAME)
		configstore.AddRemoteProvider(remoteProvider, remoteEndpoint, remotePath)
		configstore.SetConfigType("json")
		err := configstore.ReadRemoteConfig()
		if err != nil {
			log.Printf("Failed to load remote configuration from %s, Stopping application startup", remoteEndpoint)
		} else {
			appl := configstore.Get("application")
			useConfigStore = true
			log.Printf("Successfully loaded configuration from remote provider %s for %v", remoteProvider, appl)
		}
	}
	log.Printf("Loaded configuration for application %s", GetStringValue("application"))

}

func GetConfigValue(key string) interface{} {
	if useConfigStore && configstore.IsSet(key) {
		return configstore.Get(key)
	} else {
		return viper.Get(key)
	}
}

func GetStringValue(key string) string {
	if useConfigStore && configstore.IsSet(key) {
		return configstore.GetString(key)
	} else {
		return viper.GetString(key)
	}
}

func GetIntValue(key string) int {
	if useConfigStore && configstore.IsSet(key) {
		return configstore.GetInt(key)
	} else {
		return viper.GetInt(key)
	}
}

func GetFloatValue(key string) float64 {
	if useConfigStore && configstore.IsSet(key) {
		return configstore.GetFloat64(key)
	} else {
		return viper.GetFloat64(key)
	}
}

func GetBooleanValue(key string) bool {
	if useConfigStore && configstore.IsSet(key) {
		return configstore.GetBool(key)
	} else {
		return viper.GetBool(key)
	}
}

func GetStringMap(key string) map[string]interface{} {
	if useConfigStore && configstore.IsSet(key) {
		return configstore.GetStringMap(key)
	} else {
		return viper.GetStringMap(key)
	}
}
