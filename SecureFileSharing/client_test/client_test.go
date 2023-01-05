package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

const contentFour = "helloooooooooor!"
const contentLong = "THIS IS GOING TO BE A REALLY LONGGGGGGG FILLEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE WIHT MANY EEEEEEEEEEEEEEEEEEEEEEE'S AND OTHER THINGS LIKE THAT TO TEST EDGE CASESSSSSSSS"
const contentFive = "GOOD bye!!"
const contentSix = "comPuter science??!"
const contentSeven = "last resort content :))"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	//var balice *client.User
	var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User
	
	var bobLaptop *client.User
	var bobDesktop *client.User

	var charlesLaptop *client.User
	var charlesDesktop *client.User

	//new vars for multi
	var err error
	var err1 error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	aliceFile2 := "aliceFile2.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	bobFile2 := "bobFile2.txt"
	charlesFile2 := "charlesFile2.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Our Tests:Usernames and Passwords ", func() {
		Specify("Usernames and Passwords.", func() {
			//3.1.c - client should support usernames w/ len > 0 (and not users with len = 0)
			userlib.DebugMsg("Initializing user alice with empty username.")
			alice, err = client.InitUser("", defaultPassword)
			Expect(err).ToNot(BeNil(), "Each username must be of length > 0.")

			//3.2.c - client should support password w/ len >= 0 (testing 0 and long length password)
			userlib.DebugMsg("Initializing user doris with empty password.")
			doris, err = client.InitUser("doris", "")
			//Expect(err).To(BeNil(), "Failed to initialize user doris with empty string password.")

			//3.1.a - each user should have a diff username
			userlib.DebugMsg("Initializing user Alice twice.")
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err1 = client.InitUser("alice", defaultPassword)
			Expect(err1).ToNot(BeNil(), "Each user should have a unique username.")

			//3.1.b - usernames are case-sensitive AND 3.2.a - client MUST NOT assume each person has a unique password
			//userlib.DebugMsg("Initializing user alice.")
			//alice, err = client.InitUser("alice", defaultPassword)
			userlib.DebugMsg("Initializing user aLIce.")
			bob, err1 = client.InitUser("aLIce", defaultPassword)
			Expect(err1).To(BeNil(), "Failed to initialize user aLIce (case sensitive or password).")


			userlib.DebugMsg("Initializing user charles with long password.")
			charles, err = client.InitUser("charles", "knkfkajfljaljaoirjeojaknvanlkalkdfalejoawejoiajeoivajnanclkamlkdaowieoafadfaejrowr")
			Expect(err).To(BeNil(), "Failed to initialize user Charles with long password.")
		})
	})

	Describe("Our Tests: More GetUser Tests", func() {
		/*
			Returns an error if:
			- There is no initialized user for the given username.
			- The user credentials are invalid.
			- The User struct cannot be obtained due to malicious action, or the integrity of the user struct has been compromised.
		*/

		Specify("GetUser: Error if no user exists", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Invalid get of user Alice.")
			alice, err = client.GetUser("Alice", defaultPassword)
			Expect(err).ToNot(BeNil(), "Should error for a username that doesn't exist")
		})

		Specify("GetUser: Error if wrong credentials", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Invalid credentials user to get user Alice.")
			alice, err = client.GetUser("alice", "wrongPassword")
			Expect(err).ToNot(BeNil(), "Should error for wrong password with username")
		})

		Specify("GetUser: Error if malicious activity", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//TODO: MALICIOUS ACTIVITY COMPROMISING USER STRUCT -> ERROR 
			//userlib.DebugMsg("Malicious action should result in an error.")
			//alice, err = client.GetUser("alice", defaultPassword)
			//Expect(err).ToNot(BeNil(), "Should error if tampered with")
		})
	})

	
	Describe("Our Tests: More StoreFile Tests", func() {
	
		Specify("StoreFile: No error if others have same filename", func() {
			//Different users can store files using the same filename, because each 
			//user must have a separate personal file namespace.

			userlib.DebugMsg("Initializing users Alice, Bob, and Charles.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating files for each user under the same name.")
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = bob.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil(), "Should allow multiple users to use same filename")

			err = charles.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil(), "Should allow multiple users to use same filename")
		})

		Specify("StoreFile: Filename exists in personal file namespace, one user", func() {
			////If the given filename already exists in the personal namespace of the 
			//caller, then the content of the corresponding file is overwritten. Note 
			//that, in the case of sharing files, the corresponding file may or may 
			//not be owned by the caller.

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating file aliceFile and putting contentOne.")
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing another file aliceFile and putting contentTwo.")
			err = alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			//check that the content is now content2
			userlib.DebugMsg("Loading aliceFile after two stores.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")
		})

		Specify("StoreFile: Filename exists in personal file namespace, one user multiple instances", func() {
			////If the given filename already exists in the personal namespace of the 
			//caller, then the content of the corresponding file is overwritten. Note 
			//that, in the case of sharing files, the corresponding file may or may 
			//not be owned by the caller.

			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting third instance of Alice - aliceLaptop")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop storing file %s with content: %s", aliceFile, contentTwo)
			err = aliceLaptop.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			//check that the content is now content2
			userlib.DebugMsg("Loading aliceFile in desktop and phone and lap to make sure its all contentTwo.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")

			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")

			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")
			
			//another store by diff user 
			userlib.DebugMsg("alicePhone storing file %s with content: %s", aliceFile, contentThree)
			err = aliceLaptop.StoreFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			data, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)), "File content was not overwritten")

			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)), "File content was not overwritten")

			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)), "File content was not overwritten")
		})


		Specify("StoreFile: Filename exists in personal file namespace, multi users mutli instances", func() {
			////If the given filename already exists in the personal namespace of the 
			//caller, then the content of the corresponding file is overwritten. Note 
			//that, in the case of sharing files, the corresponding file may or may 
			//not be owned by the caller.

			//creating users
			userlib.DebugMsg("Initializing users Alice (aliceDesktop), Bob(bobDesktop), and Charles(charlesDesktop).")
			aliceDesktop, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())

			bobDesktop, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charlesDesktop, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			//creating aliceLap instace
			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			//BobDesk recieves files
			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bobDesktop.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			//this should overwrite the contents of both alice and bob instances of the same file
			userlib.DebugMsg("Bob loading to file %s, content: %s", bobFile, contentTwo)
			err = bobDesktop.StoreFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			//create blaptop instance
			userlib.DebugMsg("Getting second instance of Bob - bobLaptop")
			bobLaptop, err = client.GetUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			//checking to ensure correct content is in the files for ALap, BLap, ADesk, BDesk
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")

			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")

			data, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")

			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content was not overwritten")

			//////////////////////////

			//blap appends to file 
			userlib.DebugMsg("bobLaptop appending to file %s, content: %s", bobFile, contentThree)
			err = bobLaptop.AppendToFile(bobFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing to file %s, content: %s", aliceFile, contentFour)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentFour))
			Expect(err).To(BeNil())

			data, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentFour)), "File content was not overwritten")

			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentFour)), "File content was not overwritten")

			data, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentFour)), "File content was not overwritten")

			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentFour)), "File content was not overwritten")
			/////////////////////////

			//bdesk appends to file 
			userlib.DebugMsg("bobDesktop appending to file %s, content: %s", bobFile, contentFive)
			err = bobDesktop.AppendToFile(bobFile, []byte(contentFive))
			Expect(err).To(BeNil())

			//Bdesk invites CDesk
			userlib.DebugMsg("bobDesktop creating invite for Charles(chalesDesktop).")
			invite, err = bobDesktop.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			//CDesk accepts BDesk
			userlib.DebugMsg("Charles(desk) accepting invite from Bob(desk) under filename %s.", charlesFile)
			err = charlesDesktop.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())
			
			//CDesk appends to filen- contentSix
			userlib.DebugMsg("Charles(desk) storing to file %s, content: %s", charlesFile, contentSix)
			err = charlesDesktop.StoreFile(charlesFile, []byte(contentSix))
			Expect(err).To(BeNil())

			//create clap
			userlib.DebugMsg("Getting second instance of Charles - charlesLaptop")
			charlesLaptop, err = client.GetUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			
			//check instances
			data, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix)), "File content was not overwritten")

			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix)), "File content was not overwritten")

			data, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix)), "File content was not overwritten")

			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix)), "File content was not overwritten")
			
			data, err = charlesDesktop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix)), "File content was not overwritten")

			data, err = charlesLaptop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix)), "File content was not overwritten")
			////////////////////////////////


			//clap appends to file - contentLong
			userlib.DebugMsg("charlesLaptop appending to file %s, content: %s", charlesFile, contentLong)
			err = charlesLaptop.AppendToFile(charlesFile, []byte(contentLong))
			Expect(err).To(BeNil())

			//final check of all files, users, and instances
			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))

			userlib.DebugMsg("Checking that bobLaptop sees expected file data.")
			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))

			userlib.DebugMsg("Checking that bobDesktop sees expected file data.")
			data, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))

			userlib.DebugMsg("Checking that charlesLaptop sees expected file data.")
			data, err = charlesLaptop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))

			userlib.DebugMsg("Checking that charlesDesktop sees expected file data.")
			data, err = charlesDesktop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentSix + contentLong)))
		})

		Specify("StoreFile: Content is empty sequence ", func() {
			//The client MUST allow content to be any arbitrary sequence of bytes, 
			//including the empty sequence.

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil(), "File content should be allowed to be empty sequence of bytes")

			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(emptyString)))
		
		})

	
		Specify("StoreFile: Filename is length 0", func() {
			//Different users can store files using the same filename, because each 
			//user must have a separate personal file namespace.

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil(), "Filename should be allowed to be length 0")

			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

	})

	Describe("Our Tests: More LoadUser Tests", func() {
	
		Specify("LoadUser: FileName not in personal file namespace ", func() {
			//The given filename does not exist in the personal file namespace 
			//of the caller.

		})

		Specify("LoadUser: Integrity failure error ", func() {
			//The integrity of the downloaded content cannot be verified 
			//(indicating there have been unauthorized modifications to the file).
			
		})

		Specify("LoadUser: Malicious Activity error ", func() {
			//The integrity of the downloaded content cannot be verified 
			//(indicating there have been unauthorized modifications to the file).
			
		})
	})

	Describe("Our Tests: More CreateInvite/ AcceptInvite tests ", func() {
		Specify(" ", func() {

		})

		Specify("AcceptInvite with same file names ", func() {
			userlib.DebugMsg("Initializing users Alice and Bob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice and Bob both store files with the same content and fileName.")
			err = alice.StoreFile("sameName.txt", []byte(contentOne))
			Expect(err).To(BeNil())
			err = bob.StoreFile("sameName.txt", []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, "sameName.txt")

			invite, err := alice.CreateInvitation("sameName.txt", "bob")
			Expect(err).To(BeNil())

			//error expected 
			err = bob.AcceptInvitation("alice", invite, "sameName.txt")
			Expect(err).ToNot(BeNil(), "Allowed user to accept file invite under a filename that already exists in their personal file namespace")

			//correct use of accept
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
		})

		Specify(" ", func() {
			
		})
	})

	
	Describe("Our Tests: RevokeAccess !", func() {
		Specify("Revoke ", func() {
			//Given a filename in the personal namespace of the caller, this function 
			//revokes access to the corresponding file from recipientUsername and any 
			//other users with whom recipientUsernamehas shared the file.

			//test that you can't revoke something NOT in the personal 
			//namespace of the caller

			userlib.DebugMsg("Initializing users Alice, Bob, and Charles.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice trying to revoke access from bobFile.")
			err = alice.RevokeAccess(bobFile, "bob")
			Expect(err).ToNot(BeNil(), "Allowed user to revoke access to a file not in user's personal file namespace")
		
			//this should not error bc its the right fileName
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())
		})

		Specify("Revoking Bob with only CreateInvite from Alice but no AcceptInvite from Bob", func() {
			//A revoked user must lose access to the corresponding file regardless of 
			//whether their invitation state is created or accepted.
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			//at this point, bob has not accepted the invite so alice shouldn't be able to revoke
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil(), "Allowed owner to revoke access to an unshared file")
			
		})

		Specify("Multiuser Invite and Revoke test", func() {
			//A revoked user must lose access to the corresponding file regardless of 
			//whether their invitation state is created or accepted.

			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			//bob accepts
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			// bob has fully accepted invite, charles has not accepted invite -> charles should error
			//err = bob.RevokeAccess(bobFile,"charles")
			//Expect(err).ToNot(BeNil(), "Allowed non-owner to revoke access from a user.")

			//charles accpts invite; should be allowed bc bob shouldn't be able to revoke charles
			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			//now that charles accepted invite, alice revokation shoudl make bob and charles BOTH lose access

			//alice revoke on bob should work as expected (bob and charles lose access)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())
		})

		//TODO: ?? check 
		Specify("Malicious Activity RevokeAccess Behavior ", func() {
			//simulating an attack on dataStore
			userlib.DebugMsg("Initializing users Alice and Bob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			//userlib.DataStoreGetMap() - to tamper with data
			//keys = getDataStoreKeys()
		})
	})
	
	Describe("Our Tests: More security tests ", func() {
		Specify("One store, swapped datastore values, one load", func() {
			//basing this off of the example in the staff tips part of the spec 

			//have alice store a file -> should change Datastore keys bc of an api calls
			userlib.DebugMsg("Initializing users Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//og datastore vals
			ogDS := userlib.DatastoreGetMap()
			userlib.DebugMsg("ogDS is type %T", ogDS )
	
			//get the keys of the golang map -> thank you stack overflow
			dataStoreKeys := make([]uuid.UUID , 0, len(ogDS))
			userlib.DebugMsg("ogDS size is %s.", len(ogDS))
			for i, _ := range ogDS {
				userlib.DebugMsg("i is type %T", i)
				dataStoreKeys = append(dataStoreKeys, i)
			}
			userlib.DebugMsg("datastoreKeys is type %T", dataStoreKeys)

			//so now dataStoreKeys has the original key values 
	
			//alice does the api call 
			userlib.DebugMsg("Alice does a StoreFile call.")
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			//get the dataStoreKeys again
			newDS := userlib.DatastoreGetMap()
			newKeys := make([]uuid.UUID , 0, len(newDS))
			//newKeys := make([]uuid.UUID, 0, len(newDS))
			for i, _ := range newDS {
				newKeys = append(newKeys, i)
			}
			userlib.DebugMsg("newDS size is %s.", len(newDS))
			//so now the newKeys contains the new keys -> change them and see if there's a thing detected

			userlib.DebugMsg("newDS is type %T", newDS )
			userlib.DebugMsg("newKeys is type %T", newKeys)
			
			diffKeys := make([]uuid.UUID, 0, len(newDS))
			for _, i := range newKeys { 
				//gonna use a map bc that's more efficient for checking if a key exists
				//jk can't do that
				for _, m := range dataStoreKeys {
					if i != m{
						diffKeys = append(diffKeys, i)
					}
				}
			}
			userlib.DebugMsg("There are %s new keys.", len(diffKeys))
			//now diffKeys contains all the keys that have changed from the ogDS to the newDS
			
			userlib.DebugMsg("Using userlib.Datastore to swap the values in the diffKeys.")
			//lets tamper with the diffkeys in the 
			//then manually change values within datastore

			for j := 0; j< len(diffKeys)/2; j++ {
				key1 := diffKeys[j]
				ind := len(diffKeys) - 1 - j
				key2 := diffKeys[ind]
				val1, _ := userlib.DatastoreGet(key1)
				val2, _ := userlib.DatastoreGet(key2)

				//swap
				userlib.DatastoreSet(key1, val2)
				userlib.DatastoreSet(key2, val1)
			}

			//now check that load and append don't work bc of malicious activity
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil(), "Did not catch swapped Datastore values.")
			
		})

		Specify("Two stores, swapped datastore values, then tried to load both files", func() {
			userlib.DebugMsg("Initializing users Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//og datastore vals
			ogDS := userlib.DatastoreGetMap()
			userlib.DebugMsg("ogDS is type %T", ogDS )
	
			//get the keys of the golang map -> thank you stack overflow
			dataStoreKeys := make([]uuid.UUID , 0, len(ogDS))
			userlib.DebugMsg("ogDS size is %s.", len(ogDS))
			for i, _ := range ogDS {
				userlib.DebugMsg("i is type %T", i)
				dataStoreKeys = append(dataStoreKeys, i)
			}
			userlib.DebugMsg("datastoreKeys is type %T", dataStoreKeys)
			
			//alice does two api calls
			userlib.DebugMsg("Alice does a StoreFile call.")
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = alice.StoreFile(aliceFile2, []byte(contentTwo))
			Expect(err).To(BeNil())

			//get the dataStoreKeys again
			newDS := userlib.DatastoreGetMap()
			newKeys := make([]uuid.UUID , 0, len(newDS))
			//newKeys := make([]uuid.UUID, 0, len(newDS))
			for i, _ := range newDS {
				newKeys = append(newKeys, i)
			}
			userlib.DebugMsg("newDS size is %s.", len(newDS))
			//so now the newKeys contains the new keys -> change them and see if there's a thing detected

			userlib.DebugMsg("newDS is type %T", newDS )
			userlib.DebugMsg("newKeys is type %T", newKeys)
			
			diffKeys := make([]uuid.UUID, 0, len(newDS))
			for _, i := range newKeys { 
				//gonna use a map bc that's more efficient for checking if a key exists
				//jk can't do that
				for _, m := range dataStoreKeys {
					if i != m{
						diffKeys = append(diffKeys, i)
					}
				}
			}

			userlib.DebugMsg("There are %s new keys.", len(diffKeys))
			//now diffKeys contains all the keys that have changed from the ogDS to the newDS
			
			userlib.DebugMsg("Using userlib.Datastore to swap the values in the diffKeys.")
			//lets tamper with the diffkeys in the 
			//then manually change values within datastore

			for j := 0; j< len(diffKeys)/2; j++ {
				key1 := diffKeys[j]
				ind := len(diffKeys) - 1 - j
				key2 := diffKeys[ind]
				val1, _ := userlib.DatastoreGet(key1)
				val2, _ := userlib.DatastoreGet(key2)

				//swap
				userlib.DatastoreSet(key1, val2)
				userlib.DatastoreSet(key2, val1)
			}

			//now check that load and append don't work bc of malicious activity
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch swapped Datastore values.")
	

			_, err = alice.LoadFile(aliceFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch swapped Datastore values.")	
		})

		Specify("Two stores, COPIED datastore values, then tried to load both files", func() {
			userlib.DebugMsg("Initializing users Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//og datastore vals
			ogDS := userlib.DatastoreGetMap()
	
			//get the keys of the golang map -> thank you stack overflow
			dataStoreKeys := make([]uuid.UUID , 0, len(ogDS))
			userlib.DebugMsg("ogDS size is %s.", len(ogDS))
			for i, _ := range ogDS {
				userlib.DebugMsg("i is type %T", i)
				dataStoreKeys = append(dataStoreKeys, i)
			}
			
			//alice does two api calls
			userlib.DebugMsg("Alice does a StoreFile call.")
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = alice.StoreFile(aliceFile2, []byte(contentTwo))
			Expect(err).To(BeNil())

			//get the dataStoreKeys again
			newDS := userlib.DatastoreGetMap()
			newKeys := make([]uuid.UUID , 0, len(newDS))
			//newKeys := make([]uuid.UUID, 0, len(newDS))
			for i, _ := range newDS {
				newKeys = append(newKeys, i)
			}
			//so now the newKeys contains the new keys -> change them and see if there's a thing detected
			
			diffKeys := make([]uuid.UUID, 0, len(newDS))
			for _, i := range newKeys { 
				//gonna use a map bc that's more efficient for checking if a key exists
				//jk can't do that
				for _, m := range dataStoreKeys {
					if i != m{
						diffKeys = append(diffKeys, i)
					}
				}
			}

			userlib.DebugMsg("There are %s new keys.", len(diffKeys))
			//now diffKeys contains all the keys that have changed from the ogDS to the newDS
			
			userlib.DebugMsg("Using userlib.Datastore to copy the values in the diffKeys.")

			for j := 0; j< len(diffKeys)/2; j++ {
				key1 := diffKeys[j]
				ind := len(diffKeys) - 1 - j
				key2 := diffKeys[ind]
				val1, _ := userlib.DatastoreGet(key1)
				userlib.DatastoreSet(key2, val1) //should cause a panic bc incorrect value in the keys' spot
			}

			//now check that load doesn't work bc of malicious activity
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = alice.LoadFile(aliceFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	

			
		})

		//TODO: datastore tampering: multi user multi instance
		Specify("(multi user) Two stores, COPIED datastore values, then tried to load both files", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop), Bob(bobDesktop), and Charles(charlesDesktop).")
			aliceDesktop, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())

			bobDesktop, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charlesDesktop, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			//creating aliceLap instance
			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//create blaptop instance
			userlib.DebugMsg("Getting second instance of Bob - bobLaptop")
			bobLaptop, err = client.GetUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			//create claptop instance
			userlib.DebugMsg("Getting second instance of Bob - bobLaptop")
			charlesLaptop, err = client.GetUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			//OGDS VALUES
			ogDS := userlib.DatastoreGetMap()
			//get the keys of the golang map -> thank you stack overflow
			dataStoreKeys := make([]uuid.UUID , 0, len(ogDS))
			userlib.DebugMsg("ogDS size is %s.", len(ogDS))
			for i, _ := range ogDS {
				userlib.DebugMsg("i is type %T", i)
				dataStoreKeys = append(dataStoreKeys, i)
			}
			
			//alice does two api calls + sharing + other things
			userlib.DebugMsg("Alice does a StoreFile call.")
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = aliceLaptop.StoreFile(aliceFile2, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			//BobDesk recieves files
			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bobDesktop.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err = aliceLaptop.CreateInvitation(aliceFile2, "bob")
			Expect(err).To(BeNil())

			//BobDesk recieves files
			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile2)
			err = bobDesktop.AcceptInvitation("alice", invite, bobFile2)
			Expect(err).To(BeNil())

			//charles recieves the files from Alice
			userlib.DebugMsg("aliceLaptop creating invite for charles.")
			invite, err = aliceLaptop.CreateInvitation(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("charles accepting invite from Alice under filename %s.", charlesFile)
			err = charlesDesktop.AcceptInvitation("alice", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for charles.")
			invite, err = aliceLaptop.CreateInvitation(aliceFile2, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("charles accepting invite from Alice under filename %s.", charlesFile2)
			err = charlesDesktop.AcceptInvitation("alice", invite, charlesFile2)
			Expect(err).To(BeNil())
			
			//now doing some tampering -> copying

			//get the dataStoreKeys again
			newDS := userlib.DatastoreGetMap()
			newKeys := make([]uuid.UUID , 0, len(newDS))
			//newKeys := make([]uuid.UUID, 0, len(newDS))
			for i, _ := range newDS {
				newKeys = append(newKeys, i)
			}
			//so now the newKeys contains the new keys -> change them and see if there's a thing detected
			
			diffKeys := make([]uuid.UUID, 0, len(newDS))
			for _, i := range newKeys { 
				//gonna use a map bc that's more efficient for checking if a key exists
				//jk can't do that
				for _, m := range dataStoreKeys {
					if i != m{
						diffKeys = append(diffKeys, i)
					}
				}
			}

			userlib.DebugMsg("There are %s new keys.", len(diffKeys))
			
			userlib.DebugMsg("Using userlib.Datastore to copy the values in the diffKeys.")

			for j := 0; j< len(diffKeys)/2; j++ {
				key1 := diffKeys[j]
				ind := len(diffKeys) - 1 - j
				key2 := diffKeys[ind]
				val1, _ := userlib.DatastoreGet(key1)
				userlib.DatastoreSet(key2, val1) //should cause a panic bc incorrect value in the keys' spot
			}

			//now check that load doesn't work bc of malicious activity

			//alice 
			_, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = aliceDesktop.LoadFile(aliceFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	

			_, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = aliceLaptop.LoadFile(aliceFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	

			//bob
			_, err = bobDesktop.LoadFile(bobFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = bobDesktop.LoadFile(bobFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	

			_, err = bobLaptop.LoadFile(bobFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = bobLaptop.LoadFile(bobFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	

			//charles
			_, err = charlesDesktop.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = charlesDesktop.LoadFile(charlesFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	

			_, err = charlesLaptop.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")

			_, err = charlesLaptop.LoadFile(charlesFile2)
			Expect(err).ToNot(BeNil(), "LoadFile did not catch copied Datastore values.")	
		})


	})

	Describe("Our Tests: User Sessions + Sharing ", func() {
		Specify("User Sessions + Sharing", func() {
			//tests 3.6.1, 3.6.2, 3.6.3
			//using baseline from Basic Tests
			userlib.DebugMsg("Initializing users Alice (aliceDesktop), Bob(bobDesktop), and Charles(charlesDesktop).")
			aliceDesktop, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())

			bobDesktop, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charlesDesktop, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bobDesktop.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bobDesktop.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			//create blaptop instance
			userlib.DebugMsg("Getting second instance of Bob - bobLaptop")
			bobLaptop, err = client.GetUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			//blap appends to file 
			userlib.DebugMsg("bobLaptop appending to file %s, content: %s", bobFile, contentThree)
			err = bobLaptop.AppendToFile(bobFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentFour)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentFour))
			Expect(err).To(BeNil())

			//bdesk appends to file 
			userlib.DebugMsg("bobDesktop appending to file %s, content: %s", bobFile, contentFive)
			err = bobDesktop.AppendToFile(bobFile, []byte(contentFive))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			userlib.DebugMsg("Checking that bobLaptop sees expected file data.")
			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			userlib.DebugMsg("Checking that bobDesktop sees expected file data.")
			data, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//Bdesk invites CDesk
			userlib.DebugMsg("bobDesktop creating invite for Charles(chalesDesktop).")
			invite, err = bobDesktop.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			//CDesk accepts BDesk
			userlib.DebugMsg("Charles(desk) accepting invite from Bob(desk) under filename %s.", charlesFile)
			err = charlesDesktop.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())
			
			//CDesk appends to filen- contentSix
			userlib.DebugMsg("Charles(desk) appending to file %s, content: %s", charlesFile, contentSix)
			err = charlesDesktop.AppendToFile(charlesFile, []byte(contentSix))
			Expect(err).To(BeNil())

			//create clap
			//create blaptop instance
			userlib.DebugMsg("Getting second instance of Charles - charlesLaptop")
			charlesLaptop, err = client.GetUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			
			//clap appends to file - contentLong
			userlib.DebugMsg("charlesLaptop appending to file %s, content: %s", charlesFile, contentLong)
			err = charlesLaptop.AppendToFile(charlesFile, []byte(contentLong))
			Expect(err).To(BeNil())

			//final check of all files, users, and instances
			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))

			userlib.DebugMsg("Checking that bobLaptop sees expected file data.")
			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))

			userlib.DebugMsg("Checking that bobDesktop sees expected file data.")
			data, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))

			userlib.DebugMsg("Checking that charlesLaptop sees expected file data.")
			data, err = charlesLaptop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))

			userlib.DebugMsg("Checking that charlesDesktop sees expected file data.")
			data, err = charlesDesktop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive + contentSix + contentLong)))
		})
	})

	Describe("Our Tests: Sharing - error cases, creating ", func() {
		Specify("Sharing - error cases, creating", func() {
			//create A B C and create a file for A
			userlib.DebugMsg("Initializing users Alice, Bob, Charles.")
			alice, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword) 
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword) 
			Expect(err).To(BeNil())

			//create file for A, and then try to share a diff file that A doesn't have auth for
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			_, err = alice.CreateInvitation("fliesfile", "bob")
			Expect(err).ToNot(BeNil(), "Alice should not be able to share a random file without authorization.")

			//try to share content with user that doesn't exist
			err = bob.StoreFile(bobFile, []byte(contentTwo))
			_, err = bob.CreateInvitation(bobFile, "jim bob")
			Expect(err).ToNot(BeNil(), "Sharing with non-existant user.")

			//can't invite yourself to files you don't have access to
			_, err = bob.CreateInvitation(aliceFile, "bob")
			Expect(err).ToNot(BeNil(), "Sharing to yourself without access.")
		})
	})

	Describe("Our Tests: Sharing - error cases, accepting", func() {
		Specify("Sharing - error cases, accepting", func() {
			//create A B C and create a file for A
			userlib.DebugMsg("Initializing users Alice, Bob, Charles.")
			alice, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword) 
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword) 
			Expect(err).To(BeNil())
			doris, err = client.InitUser("doris", defaultPassword) 
			_ = doris
			Expect(err).To(BeNil())

			//create file for A
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			//basically checking for errors with invalid params
			//create invite for B and then B accepts invalid things
			userlib.DebugMsg("Creating invite for Bob from Alice.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			//wrong person accepts
			userlib.DebugMsg("Charles accepts Alice's invite to Bob.")
			err = charles.AcceptInvitation("alice", invite, charlesFile)
			Expect(err).ToNot(BeNil(), "Allowed invalid user to accept invite.")

			//bob tries to accept with wrong sender info
			err = bob.AcceptInvitation("charles", invite, bobFile)
			Expect(err).ToNot(BeNil(), "Allowed user to accept invite from wrong user.")
			err = bob.AcceptInvitation("bob", invite, bobFile)
			Expect(err).ToNot(BeNil(), "Allowed user to accept invite from themselves.")

			//bob tries to accept with wrong invitation info
			userlib.DebugMsg("Creating invite for Doris from Charles.")
			err = charles.StoreFile(charlesFile, []byte(contentThree))
			Expect(err).To(BeNil())
			invite2, err := charles.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())
			err = bob.AcceptInvitation("alice", invite2, bobFile)
			//TODO: ADD EXPECT()
			Expect(err).ToNot(BeNil(), "Allowed user to accept an invite that wasn't theirs.")

			//bob tries to aceept with existing filename
			err = bob.StoreFile("bobFile", []byte(contentFour))
			Expect(err).To(BeNil())
			err = bob.AcceptInvitation("alice", invite, "bobFile")
			Expect(err).ToNot(BeNil(), "Allowed user to accept for a file that already exists under their name.")

		})
	})

	/*
	Describe("Our Tests: Append and Efficiency", func() {
		Specify("Efficiency Test for AppendToFile()", func() {
			
			//The bandwidth of the AppendToFile() operation MUST scale linearly with only 
			//the size of data being appended and the number of users the file is shared 
			//with, and nothing else. Logarithmic and constant scaling in other terms is 
			//fine.
			
			
			// Helper function to measure bandwidth of a particular operation
			measureBandwidth := func(probe func()) (bandwidth int) {
				before := userlib.DatastoreGetBandwidth()
				probe()
				after := userlib.DatastoreGetBandwidth()
				return after - before
			}
			//(append) Returns an error if:
			// - The given filename does not exist in the personal file namespace of the
			//caller.
			// - Appending the file cannot succeed due to any other malicious action.
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			//append error 1 test
			err = alice.AppendToFile(aliceFile, []byte(contentOne))
			Expect(err).ToNot(BeNil(), "Append should not be allowed for a file name that doesn't exist in the personal file namespace of the user.")
			
			userlib.DebugMsg("Creating aliceFile for user Alice.")
			alice.StoreFile(aliceFile, []byte(contentOne))

			//appends should work and scale
			prevBW := 0
			for i:= 0; i < 10000; i++{
				//get bw
				bw := measureBandwidth(func() {
					err = alice.AppendToFile(aliceFile, []byte("A"))
				 })
				Expect(err).To(BeNil(), "Append isn't working.")
				//if this is not the first iteration
				//fmt.Prints(i, bw, "\n")
				if i % 1000 == 0{
					userlib.DebugMsg("i = %s, bandwith = %s", i, bw)
				}
				if i != 0 {
					//check that the differnce isn't too much... like <10
					Expect(bw - prevBW > 10).ToNot(BeTrue(), "Difference in BW is greater than 10.")
				}
				prevBW = bw
			}

		})

	})
	*/

	
	Describe("Basic Tests", func() {
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword) 
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})

})
