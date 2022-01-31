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
// SubscriberModel provide firestore data model to handle Subscriber specific CRUD operations.

package models

import (
	"context"
	"encoding/json"
	"log"
)

type Subscriber struct {
	Name         string `json:"name" csv:"Name" firestore:"name"`
	Email        string `json:"email" csv:"Subscribers" firestore:"email"`
	CreationDate string `json:"creationdate" csv:"-" firestore:"creationdate"`
	IsActive     bool   `json:"status" csv:"-" firestore:"status"`
}

const SUBSCRIBER_DATA_COLLECTION string = "subscribers"

func NewSubscriber(email string) *Subscriber {
	subscriber := Subscriber{Email: email, IsActive: true}
	return &subscriber
}

func SaveSubscriberData(data interface{}, id string) bool {
	var dataMap map[string]interface{}
	dj, _err := json.Marshal(data)
	if _err != nil {
		log.Println("FATAL: JSON data marshalling failed, SubscriberModel.SaveSubscriberData")
		panic(_err)
	}
	json.Unmarshal(dj, &dataMap)
	ctx := context.Background()
	dbClient := GetFirestoreClient(ctx, GetCloudProject())
	_status := saveOrUpdateEntityWithId(ctx, dbClient, SUBSCRIBER_DATA_COLLECTION, id, dataMap)
	defer dbClient.Close()
	if _status == RECORD_SAVED {
		return true
	} else {
		return false
	}
}

func getSubscriber(emailId string) *Subscriber {
	var subscriberData Subscriber
	ctx := context.Background()
	dbClient := GetFirestoreClient(ctx, GetCloudProject())
	docsn := getEntityWithId(ctx, dbClient, SUBSCRIBER_DATA_COLLECTION, emailId)
	defer dbClient.Close()
	if docsn != nil {
		docsn.DataTo(&subscriberData)
		return &subscriberData
	}
	return nil
}

func FindAllActiveSubscribers() []Subscriber {
	var subscriberData Subscriber
	var subscribers []Subscriber
	ctx := context.Background()
	dbClient := GetFirestoreClient(ctx, GetCloudProject())
	_docsns := findAllMatchingEntity(ctx, dbClient, SUBSCRIBER_DATA_COLLECTION, "status", "==", true)
	defer dbClient.Close()
	if _docsns == nil {
		return nil
	}
	for _, doc := range _docsns {
		doc.DataTo(&subscriberData)
		subscribers = append(subscribers, subscriberData)
	}
	return subscribers
}
