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

// FirestoreModel initialize firestore client and provide generic methods for CRUD operations.

package models

import (
	"ail-examples/config"
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	RECORD_INSERTED  string = "success"
	RECORD_SAVED     string = "success"
	RECORD_DELETED   string = "success"
	ERROR_INSERT_OPS string = "RECORD-INSERT-FAILED"
	ERROR_SAVE_OPS   string = "RECORD-SAVE-FAILED"
	ERROR_QUERY_OPS  string = "RECORD-QUERY-FAILED"
	ERROR_GET_OPS    string = "RECORD-READ-FAILED"
	ERROR_DELETE_OPS string = "RECORD-DELETE-FAILED"
	QUERY_MAX_LIMIT  int    = 100
)

var projectID string

func DBInit() {
	// Use GCP service account
	ctx := context.Background()
	projectID := GetCloudProject()
	dataStoreClient := GetFirestoreClient(ctx, projectID)
	defer dataStoreClient.Close()
}

func GetCloudProject() string {
	if projectID != "" {
		return projectID
	}
	projectID = config.GetStringValue(config.KEY_PROJECT_ID)
	if projectID == "" {
		log.Fatalf("FATAL: Required configuration property not found: %s", config.KEY_PROJECT_ID)
	}
	return projectID
}

func GetFirestoreClient(cloudContext context.Context, project string) *firestore.Client {
	dbClient, err := firestore.NewClient(cloudContext, project)
	if err != nil {
		log.Println("FATAL: Failed to create firestore client")
		log.Fatalln(err)
	}
	return dbClient
}

func saveOrUpdateEntityWithId(ctx context.Context, client *firestore.Client, entityName string, id string, data map[string]interface{}) string {
	var msg = RECORD_SAVED
	_success, err := client.Collection(entityName).Doc(id).Set(ctx, data, firestore.MergeAll)
	if err != nil {
		msg = ERROR_SAVE_OPS
		log.Printf("ERROR: %s %v", msg, err)
	} else {
		log.Println("DEBUG: Entity with Id " + id + " saved : " + _success.UpdateTime.String())
	}
	return msg
}

func findAllMatchingEntity(ctx context.Context, client *firestore.Client, entityName string, queryField string, operation string, value interface{}) []*firestore.DocumentSnapshot {
	docsn, err := client.Collection(entityName).Where(queryField, operation, value).Limit(QUERY_MAX_LIMIT).Documents(ctx).GetAll()
	if err != nil {
		if status.Code(err) != codes.NotFound {
			log.Printf("ERROR: %s %v", ERROR_QUERY_OPS, err)
		}
		// No record found
		log.Printf("DEBUG: No record found for the query: %s where %s %s %s ", entityName, queryField, operation, value)
		return nil
	}
	return docsn
}

func getEntityWithId(ctx context.Context, client *firestore.Client, entityName string, id string) *firestore.DocumentSnapshot {
	doc, err := client.Collection(entityName).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			// log error except record not found
			log.Printf("ERROR: %s %v", ERROR_GET_OPS, err)
		}
		return nil
	}
	return doc
}

func DoMoneyTransaction(ctx context.Context, client *firestore.Client,
	senderAccountRef *firestore.DocumentRef,
	receiverAccountRef *firestore.DocumentRef, amount float64) error {

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		senderAcctDoc, serr := tx.Get(senderAccountRef)
		receiverAcctDoc, rerr := tx.Get(receiverAccountRef)
		if serr != nil || rerr != nil {
			return serr
		}
		sbal, errs := senderAcctDoc.DataAt("balance")
		rbal, errr := receiverAcctDoc.DataAt("balance")
		if errs != nil || errr != nil {
			return errs
		}
		senderBal := sbal.(float64)
		receiverBal := rbal.(float64)
		if senderBal > amount {
			serr = tx.Set(senderAccountRef, map[string]interface{}{
				"balance": senderBal - amount,
			}, firestore.MergeAll)
			rerr = tx.Set(receiverAccountRef, map[string]interface{}{
				"balance": receiverBal + amount,
			}, firestore.MergeAll)
			if serr == nil {
				return rerr
			}
			return serr
		}
		return errors.New("Insufficient balance in sender account")
	})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("Transaction error has occurred: %s", err)
	}
	return err
}
