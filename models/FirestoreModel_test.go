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

package models

import (
	"ail-examples/config"
	"context"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

const TEST_COLL_NAME = "testaccount"

func init() {
	viper.AddConfigPath("..")
	viper.AddConfigPath("../resources")
	config.AppInit()
	// Initialize Database Configuration
	DBInit()
}

func TestDoMoneyTransaction(t *testing.T) {
	sender := map[string]interface{}{
		"id":            "1001",
		"bankcode":      "AWGP",
		"accountnumber": "1001000401",
		"balance":       1300000.00,
		"name":          "Om Prakash",
	}
	receiver := map[string]interface{}{
		"id":            "1013",
		"bankcode":      "UNIV",
		"accountnumber": "3100100013",
		"balance":       130000000.00,
		"name":          "Aadhyarupam",
	}
	ctx := context.Background()
	dbClient := GetFirestoreClient(ctx, GetCloudProject())
	senderId := sender["id"].(string)
	receiverId := receiver["id"].(string)
	_status := saveOrUpdateEntityWithId(ctx, dbClient, TEST_COLL_NAME, senderId, sender)
	_rstatus := saveOrUpdateEntityWithId(ctx, dbClient, TEST_COLL_NAME, receiverId, receiver)
	if _status == RECORD_SAVED && _rstatus == RECORD_SAVED {
		sdoc := dbClient.Collection(TEST_COLL_NAME).Doc(senderId)
		rdoc := dbClient.Collection(TEST_COLL_NAME).Doc(receiverId)
		// Test PASS
		err := DoMoneyTransaction(ctx, dbClient, sdoc, rdoc, 10000)
		if err != nil {
			t.Error(err)
		}

		// Test "Insufficient balance case"
		err = DoMoneyTransaction(ctx, dbClient, sdoc, rdoc, 20000000)
		if err == nil {
			t.Error("Expected error for insufficent balance in sender account")
		} else {
			fmt.Printf("Error occured as expected, %v", err)
		}

	} else {
		t.Error("Failed to save account data")
		return
	}

}
