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
package main

import (
	"ail-examples/config"
	"ail-examples/models"
	"log"
	"net/http"
)

func main() {
	initialize()
	mux := http.NewServeMux()
	log.Println("Starting HTTP subscription service ...")
	mux.HandleFunc("/api/subscribe", subscriptionHandler)
	log.Fatal(http.ListenAndServe(":8010", mux))
}

func initialize() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	config.AppInit()
	// Initialize Database Configuration
	models.DBInit()
	// Initialize Mail Server
	initMailServer()
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation
	log.Println("subscriptionHandler: Request received ...")
	query := r.URL.Query()
	email := query.Get("email")
	name := query.Get("name")

	if !isValidEmail(email) {
		msg := "Please provide valid Email"
		log.Printf("subscriptionHandler: %s", msg)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}
	if createEmailSubscription(email, name) {
		msg := "Subscription created successfully for Email " + email
		log.Printf("subscriptionHandler: %s", msg)
		w.Write([]byte(msg))
	}
}

func isValidEmail(e string) bool {
	if len(e) < 5 || len(e) > 254 {
		return false
	}
	return true
}

func isValidName(name string) bool {
	if len(name) < 2 || len(name) > 50 {
		return false
	}
	return true
}

func createEmailSubscription(email string, name string) bool {
	// Create subscription
	subscriber := models.NewSubscriber(email)
	if isValidName(name) {
		subscriber.Name = name
	}
	result := models.SaveSubscriberData(subscriber, email)
	return result
}

func initMailServer() {
	log.Println("Initialized Mail server connection ...")
	// TBD
}
