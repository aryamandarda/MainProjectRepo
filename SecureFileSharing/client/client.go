package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	//"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	ID uuid.UUID
	Username string
	Password []byte
	Salt []byte
	Pkey userlib.PKEDecKey
	DSkey userlib.DSSignKey

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

// NOTE: The following methods have toy (insecure!) implementations.
//NOTE: UUIDs are linked lists. HASHKDF used to spread keys from root key to shared users. Root key stored with creator.
func InitUser(username string, password string) (userdataptr *User, err error) {
	if len(username) < 1 {
		return nil, errors.New("Cannot provide an empty string as a username!")
	}
	_, exists := userlib.KeystoreGet(username)
	if exists {
		return nil, errors.New("The username: '" + username + "' already exists.")
	}
	var userdata User
	userdata.Username = username
	userdata.Salt = userlib.Hash([]byte(username))
	userdata.Password = userlib.Argon2Key([]byte(password), userdata.Salt, 32)
	key := userdata.Password[:16]
	mac := userdata.Password[16:]
	pubKey, privKey, err := userlib.PKEKeyGen()
	if err != nil {
		return nil, err
	}
	privSign, pubSign, err := userlib.DSKeyGen()
	if err != nil {
		return nil, err
	}
	userdata.Pkey = privKey
	userdata.DSkey = privSign
	err = userlib.KeystoreSet(username + "0", pubKey)
	if err != nil {
		return nil, err
	}
	err = userlib.KeystoreSet(username + "1", pubSign)
	if err != nil {
		return nil, err
	}
	user_UUID, err := uuid.FromBytes([]byte(username + password)[:16])
	if err != nil {
		return nil, err
	}
	userdata.ID = user_UUID

	// Storing info in DataStore
	marshalled_userData, err := json.Marshal(&userdata)
	if err != nil {
		return nil, err
	}
	random_iv := userlib.RandomBytes(16)
	encryted_data := userlib.SymEnc(key, random_iv, marshalled_userData)
	hmac_data, err := userlib.HMACEval(mac, encryted_data)
	if err != nil {
		return nil, err
	}
	storage := append(hmac_data, encryted_data...)
	userlib.DatastoreSet(user_UUID, storage)
	
	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata
	user_UUID, err := uuid.FromBytes([]byte(username + password)[:16])
	if err != nil {
		return nil, err
	}
	fetched, ok := userlib.DatastoreGet(user_UUID)
	if !ok {
		return nil, errors.New("No user found. Username or password is incorrect or user does not exist.")
	}
	salt := userlib.Hash([]byte(username))
	pw := userlib.Argon2Key([]byte(password), salt, 32)
	key := pw[:16]
	mac := pw[16:]
	if len(fetched) < 65 {
		return nil, errors.New("Data dont make sense bruh.")
	}
	hmac, err := userlib.HMACEval(mac, fetched[64:])
	if err != nil {
		return nil, err
	}
	equal := userlib.HMACEqual(fetched[:64], hmac)
	if !equal {
		return nil, errors.New("Tampering of data detected! Cannot retrieve user.")
	}
	encrypted_data := fetched[64:]
	decrypted := userlib.SymDec(key, encrypted_data)
	err = json.Unmarshal(decrypted, userdataptr)
	if err != nil {
		return nil, err
	}
	return userdataptr, nil
}

// Useful structs for files
type File struct {
	Owner string
	ID uuid.UUID
	Parts int
	Next uuid.UUID
	Last uuid.UUID
	LastFileNode uuid.UUID
	UserFileMap map[string][]byte
	SharedUserMap map[string][]string
}

type FileAccess struct {
	FileLoc uuid.UUID
	Key []byte
	MAC []byte
}

type FilePart struct {
	Content []byte
	Next uuid.UUID
}

// HELPER FUNCTIONS
//
// ENCRYPTION
func EncryptFilePart(fileAccess FileAccess, file File, filePart FilePart, userdata *User) (encryptedFilePart []byte, err error) {
	marshaled_filePart, err := json.Marshal(&filePart)
	if err != nil {
		return nil, err
	}
	derivedHashKDF, err := userlib.HashKDF(fileAccess.Key, []byte("Create chain of file nodes encryption"))
	if err != nil {
		return nil, err
	}
	partKey := derivedHashKDF[:16]
	partMAC := derivedHashKDF[16:32]
	randomFileIV := userlib.RandomBytes(16)
	encrytedFilePartData := userlib.SymEnc(partKey, randomFileIV, marshaled_filePart)
	hmacFilePartData, err := userlib.HMACEval(partMAC, encrytedFilePartData)
	if err != nil {
		return nil, err
	}
	filePartStorage := append(hmacFilePartData, encrytedFilePartData...)
	return filePartStorage, nil
}

func EncryptFile(fileAccess FileAccess, file File, userdata *User) (encryptedFile []byte, err error) {
	marshaled_file, err := json.Marshal(&file)
	if err != nil {
		return nil, err
	}
	randomFileIV := userlib.RandomBytes(16)
	encrytedFileData := userlib.SymEnc(fileAccess.Key, randomFileIV, marshaled_file)
	hmacFileData, err := userlib.HMACEval(fileAccess.MAC, encrytedFileData)
	if err != nil {
		return nil, err
	}
	fileStorage := append(hmacFileData, encrytedFileData...)
	return fileStorage, nil
}

func EncryptFileAccess(fileAccess FileAccess, userdata *User, filename string) (encryptedFileAccess []byte, err error) {
	marshaled_fileAcess, err := json.Marshal(&fileAccess)
	if err != nil {
		return nil, err
	}
	salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key := deterministicK[:16]
	mac := deterministicK[16:]
	randomFileIV := userlib.RandomBytes(16)
	encryptedFileAccessData := userlib.SymEnc(key, randomFileIV, marshaled_fileAcess)
	hmacFileAccess, err := userlib.HMACEval(mac, encryptedFileAccessData)
	if err != nil {
		return nil, err
	}
	fileAccessStorage := append(hmacFileAccess, encryptedFileAccessData...)
	return fileAccessStorage, nil
}

// DECRYPTION

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	fileAccessID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return err
	}
	encryptedFileAccess, ok := userlib.DatastoreGet(fileAccessID)
	if !ok {
		var file File
		var filePart FilePart
		var fileAccess FileAccess

		// Populating fileAccess struct
		fileAccess.FileLoc = uuid.New()
		fileAccess.Key = userlib.RandomBytes(16)
		fileAccess.MAC = userlib.RandomBytes(16)

		// Populating file struct
		file.Owner = userdata.Username
		file.ID = fileAccess.FileLoc
		file.Parts = 1
		file.Next = uuid.New()
		file.Last = uuid.New()
		file.LastFileNode = file.Next
		file.UserFileMap = make(map[string][]byte)
		file.SharedUserMap = make(map[string][]string)

		// Populating filePart struct
		filePart.Content = content
		filePart.Next = file.Last

		// Encrypting filePart struct
		marshaled_filePart, err := json.Marshal(&filePart)
		if err != nil {
			return err
		}
		derivedHashKDF, err := userlib.HashKDF(fileAccess.Key, []byte("Create chain of file nodes encryption"))
		if err != nil {
			return err
		}
		partKey := derivedHashKDF[:16]
		partMAC := derivedHashKDF[16:32]
		randomFileIV := userlib.RandomBytes(16)
		encrytedFilePartData := userlib.SymEnc(partKey, randomFileIV, marshaled_filePart)
		hmacFilePartData, err := userlib.HMACEval(partMAC, encrytedFilePartData)
		if err != nil {
			return err
		}
		filePartStorage := append(hmacFilePartData, encrytedFilePartData...)
		userlib.DatastoreSet(file.Next, filePartStorage)

		// Encrypting file struct
		marshaled_file, err := json.Marshal(&file)
		if err != nil {
			return err
		}
		randomFileIV = userlib.RandomBytes(16)
		encrytedFileData := userlib.SymEnc(fileAccess.Key, randomFileIV, marshaled_file)
		hmacFileData, err := userlib.HMACEval(fileAccess.MAC, encrytedFileData)
		if err != nil {
			return err
		}
		fileStorage := append(hmacFileData, encrytedFileData...)
		userlib.DatastoreSet(file.ID, fileStorage)

		// Encrypting fileAccess struct
		marshaled_fileAcess, err := json.Marshal(&fileAccess)
		if err != nil {
			return err
		}
		salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
		deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
		key := deterministicK[:16]
		mac := deterministicK[16:]
		randomFileIV = userlib.RandomBytes(16)
		encryptedFileAccessData := userlib.SymEnc(key, randomFileIV, marshaled_fileAcess)
		hmacFileAccess, err := userlib.HMACEval(mac, encryptedFileAccessData)
		if err != nil {
			return err
		}
		fileAccessStorage := append(hmacFileAccess, encryptedFileAccessData...)
		userlib.DatastoreSet(fileAccessID, fileAccessStorage)
	} else {
		var fileAccessptr FileAccess
		var fileptr File
		var filePart FilePart
		// Decrypting fileAccess struct
		salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
		deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
		key := deterministicK[:16]
		mac := deterministicK[16:]
		if len(encryptedFileAccess) < 65 {
			return errors.New("Data dont make sense bruh.")
		}
		hmac, err := userlib.HMACEval(mac, encryptedFileAccess[64:])
		if err != nil {
			return err
		}
		equal := userlib.HMACEqual(hmac, encryptedFileAccess[:64])
		if !equal {
			return errors.New("Tampering of data detected! Cannot retrieve file access.")
		}
		encryptedData := encryptedFileAccess[64:]
		decrypted := userlib.SymDec(key, encryptedData)
		err = json.Unmarshal(decrypted, &fileAccessptr)
		if err != nil {
			return err
		}

		// Decrypting file struct
		encryptedFile, ok := userlib.DatastoreGet(fileAccessptr.FileLoc)
		if !ok {
			return errors.New("Could not locate file info.")
		}
		if len(encryptedFile) < 65 {
			return errors.New("Data dont make sense bruh.")
		}
		fileHMAC, err := userlib.HMACEval(fileAccessptr.MAC, encryptedFile[64:])
		if err != nil {
			return err
		}
		equal = userlib.HMACEqual(fileHMAC, encryptedFile[:64])
		if !equal {
			return errors.New("Tampering of data detected! Cannot retrieve file.")
		}
		encryptedData = encryptedFile[64:]
		decrypted = userlib.SymDec(fileAccessptr.Key, encryptedData)
		err = json.Unmarshal(decrypted, &fileptr)
		if err != nil {
			return err
		}

		// Delete all previous filePart nodes
		partKey := fileAccessptr.Key
		var partMAC []byte
		currentUUID := fileptr.Next
		var filePartptr FilePart
		for partCounter := 0; partCounter < fileptr.Parts; partCounter++ {
			derivedHashKDF, err := userlib.HashKDF(partKey, []byte("Create chain of file nodes encryption"))
			if err != nil {
				return err
			}
			partKey = derivedHashKDF[:16]
			partMAC = derivedHashKDF[16:32]
			// Decrypt fileNode
			encryptedNode, ok := userlib.DatastoreGet(currentUUID)
			if !ok {
				return errors.New(fmt.Sprintf("Could not locate file info for node %v.", partCounter))
			}
			if len(encryptedNode) < 65 {
				return errors.New("Data dont make sense bruh.")
			}
			nodeHMAC, err := userlib.HMACEval(partMAC, encryptedNode[64:])
			if err != nil {
				return err
			}
			equal = userlib.HMACEqual(nodeHMAC, encryptedNode[:64])
			if !equal {
				return errors.New("Tampering of data detected! Cannot delete filePart.")
			}
			encryptedData = encryptedNode[64:]
			decrypted = userlib.SymDec(partKey, encryptedData)
			err = json.Unmarshal(decrypted, &filePartptr)
			if err != nil {
				return err
			}
			userlib.DatastoreDelete(currentUUID)
			currentUUID = filePartptr.Next
		}

		// Populate filePart struct
		filePart.Content = content
		filePart.Next = fileptr.Last

		// Encrypt filePart
		marshaled_filePart, err := json.Marshal(&filePart)
		if err != nil {
			return err
		}
		derivedHashKDF, err := userlib.HashKDF(fileAccessptr.Key, []byte("Create chain of file nodes encryption"))
		if err != nil {
			return err
		}
		partKey = derivedHashKDF[:16]
		partMAC = derivedHashKDF[16:32]
		randomFileIV := userlib.RandomBytes(16)
		encrytedFilePartData := userlib.SymEnc(partKey, randomFileIV, marshaled_filePart)
		hmacFilePartData, err := userlib.HMACEval(partMAC, encrytedFilePartData)
		if err != nil {
			return err
		}
		filePartStorage := append(hmacFilePartData, encrytedFilePartData...)
		userlib.DatastoreSet(fileptr.Next, filePartStorage)

		// Updating file struct data
		fileptr.Parts = 1

		// Encrypting file struct
		marshaled_file, err := json.Marshal(fileptr)
		if err != nil {
			return err
		}
		randomFileIV = userlib.RandomBytes(16)
		encrytedFileData := userlib.SymEnc(fileAccessptr.Key, randomFileIV, marshaled_file)
		hmacFileData, err := userlib.HMACEval(fileAccessptr.MAC, encrytedFileData)
		if err != nil {
			return err
		}
		fileStorage := append(hmacFileData, encrytedFileData...)
		userlib.DatastoreSet(fileptr.ID, fileStorage)
	}

	return nil
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	var fileAccessptr FileAccess
	var fileptr File
	//var filePartptr *FilePart
	var filePart FilePart
	// Decrypting FileAccess struct
	ID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return err
	}
	encryptedFileAccess, ok := userlib.DatastoreGet(ID)
	if !ok {
		return errors.New("Failed to fetch from DataStore. File does not exist.")
	}
	salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key := deterministicK[:16]
	mac := deterministicK[16:]
	if len(encryptedFileAccess) < 65 {
		return errors.New("Data dont make sense bruh.")
	}
	hmac, err := userlib.HMACEval(mac, encryptedFileAccess[64:])
	if err != nil {
		return err
	}
	equal := userlib.HMACEqual(hmac, encryptedFileAccess[:64])
	if !equal {
		return errors.New("Tampering of data detected! Cannot retrieve file access.")
	}
	encryptedData := encryptedFileAccess[64:]
	decrypted := userlib.SymDec(key, encryptedData)
	err = json.Unmarshal(decrypted, &fileAccessptr)
	if err != nil {
		return err
	}

	// Decrypting file struct
	encryptedFile, ok := userlib.DatastoreGet(fileAccessptr.FileLoc)
	if !ok {
		return errors.New("Could not locate file info.")
	}
	if len(encryptedFile) < 65 {
		return errors.New("Data dont make sense bruh.")
	}
	fileHMAC, err := userlib.HMACEval(fileAccessptr.MAC, encryptedFile[64:])
	if err != nil {
		return err
	}
	equal = userlib.HMACEqual(fileHMAC, encryptedFile[:64])
	if !equal {
		return errors.New("Tampering of data detected! Cannot retrieve file.")
	}
	encryptedData = encryptedFile[64:]
	decrypted = userlib.SymDec(fileAccessptr.Key, encryptedData)
	err = json.Unmarshal(decrypted, &fileptr)
	if err != nil {
		return err
	}

	// Updating #fileParts
	fileptr.Parts = fileptr.Parts + 1

	// Calculating key for new filePart node using HashKDF
	partKey := fileAccessptr.Key
	var partMAC []byte
	for partCounter := 0; partCounter < fileptr.Parts; partCounter++ {
		derivedHashKDF, err := userlib.HashKDF(partKey, []byte("Create chain of file nodes encryption"))
		if err != nil {
			return err
		}
		partKey = derivedHashKDF[:16]
		partMAC = derivedHashKDF[16:32]
	}

	// Populating filePart struct
	filePart.Content = content
	filePart.Next = uuid.New()

	// Changing pointers in the fileParts and file list
	fileptr.LastFileNode = fileptr.Last
	fileptr.Last = filePart.Next
	
	// Encrypting filePart struct
	marshaled_filePart, err := json.Marshal(&filePart)
	if err != nil {
		return err
	}
	randomFileIV := userlib.RandomBytes(16)
	encrytedFilePartData := userlib.SymEnc(partKey, randomFileIV, marshaled_filePart)
	hmacFilePartData, err := userlib.HMACEval(partMAC, encrytedFilePartData)
	if err != nil {
		return err
	}
	filePartStorage := append(hmacFilePartData, encrytedFilePartData...)
	userlib.DatastoreSet(fileptr.LastFileNode, filePartStorage)

	// Encrypting file struct
	marshaled_file, err := json.Marshal(fileptr)
	if err != nil {
		return err
	}
	randomFileIV = userlib.RandomBytes(16)
	encrytedFileData := userlib.SymEnc(fileAccessptr.Key, randomFileIV, marshaled_file)
	hmacFileData, err := userlib.HMACEval(fileAccessptr.MAC, encrytedFileData)
	if err != nil {
		return err
	}
	fileStorage := append(hmacFileData, encrytedFileData...)
	userlib.DatastoreSet(fileptr.ID, fileStorage)

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	var fileAccessptr FileAccess
	var fileptr File
	var filePartptr FilePart

	// Decrypting FileAccess struct
	ID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return nil, err
	}
	encryptedFileAccess, ok := userlib.DatastoreGet(ID)
	if !ok {
		return nil, errors.New("Failed to fetch from DataStore. File does not exist.")
	}
	salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key := deterministicK[:16]
	mac := deterministicK[16:]
	if len(encryptedFileAccess) < 65 {
		return nil, errors.New("Data dont make sense bruh.")
	}
	if len(encryptedFileAccess) < 65 {
		return nil, errors.New("Data dont make sense bruh.")
	}
	hmac, err := userlib.HMACEval(mac, encryptedFileAccess[64:])
	if err != nil {
		return nil, err
	}
	equal := userlib.HMACEqual(hmac, encryptedFileAccess[:64])
	if !equal {
		return nil, errors.New("Tampering of data detected! Cannot retrieve file access.")
	}
	encryptedData := encryptedFileAccess[64:]
	decrypted := userlib.SymDec(key, encryptedData)
	err = json.Unmarshal(decrypted, &fileAccessptr)
	if err != nil {
		return nil, err
	}

	// Decrypting file struct
	encryptedFile, ok := userlib.DatastoreGet(fileAccessptr.FileLoc)
	if !ok {
		return nil, errors.New("Could not locate file info.")
	}
	if len(encryptedFile) < 65 {
		return nil, errors.New("Data dont make sense bruh.")
	}
	fileHMAC, err := userlib.HMACEval(fileAccessptr.MAC, encryptedFile[64:])
	if err != nil {
		return nil, err
	}
	equal = userlib.HMACEqual(fileHMAC, encryptedFile[:64])
	if !equal {
		return nil, errors.New("Tampering of data detected! Cannot retrieve file.")
	}
	encryptedData = encryptedFile[64:]
	decrypted = userlib.SymDec(fileAccessptr.Key, encryptedData)
	err = json.Unmarshal(decrypted, &fileptr)
	if err != nil {
		return nil, err
	}

	// Calculating key for first filePart node using HashKDF
	partKey := fileAccessptr.Key
	var partMAC []byte
	currentUUID := fileptr.Next
	for partCounter := 0; partCounter < fileptr.Parts; partCounter++ {
		derivedHashKDF, err := userlib.HashKDF(partKey, []byte("Create chain of file nodes encryption"))
		if err != nil {
			return nil, err
		}
		partKey = derivedHashKDF[:16]
		partMAC = derivedHashKDF[16:32]
		// Decrypt fileNode
		encryptedNode, ok := userlib.DatastoreGet(currentUUID)
		if !ok {
			return nil, errors.New("Could not locate file info.")
		}
		if len(encryptedNode) < 65 {
			return nil, errors.New("Data dont make sense bruh.")
		}
		nodeHMAC, err := userlib.HMACEval(partMAC, encryptedNode[64:])
		if err != nil {
			return nil, err
		}
		equal = userlib.HMACEqual(nodeHMAC, encryptedNode[:64])
		if !equal {
			return nil, errors.New("Tampering of data detected! Cannot retrieve filePart.")
		}
		encryptedData = encryptedNode[64:]
		decrypted = userlib.SymDec(partKey, encryptedData)
		err = json.Unmarshal(decrypted, &filePartptr)
		if err != nil {
			return nil, err
		}
		content = append(content, filePartptr.Content...)
		currentUUID = filePartptr.Next
	}

	return content, err
}

type Invite struct {
	EncryptedKey []byte
	SignedMessage []byte
	EncryptedFileAccess []byte
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	var invite Invite
	var fileAccess FileAccess
	var ownerFileAccess FileAccess
	var file File

	if userdata.Username == recipientUsername {
		return invitationPtr, errors.New("Cannot send invite to oneself.")
	}

	// Decrypting Owner FileAccess struct
	ID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return invitationPtr, err
	}
	encryptedFileAccess, ok := userlib.DatastoreGet(ID)
	if !ok {
		return invitationPtr, errors.New("Failed to fetch from DataStore. File does not exist.")
	}
	salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key := deterministicK[:16]
	mac := deterministicK[16:]
	if len(encryptedFileAccess) < 65 {
		return invitationPtr, errors.New("Data dont make sense bruh.")
	}
	hmac, err := userlib.HMACEval(mac, encryptedFileAccess[64:])
	if err != nil {
		return invitationPtr, err
	}
	equal := userlib.HMACEqual(hmac, encryptedFileAccess[:64])
	if !equal {
		return invitationPtr, errors.New("Tampering of data detected! Cannot retrieve file access.")
	}
	encryptedData := encryptedFileAccess[64:]
	decrypted := userlib.SymDec(key, encryptedData)
	err = json.Unmarshal(decrypted, &ownerFileAccess)
	if err != nil {
		return invitationPtr, err
	}

	// Decrypting file struct
	encryptedFile, ok := userlib.DatastoreGet(ownerFileAccess.FileLoc)
	if !ok {
		return invitationPtr, errors.New("Could not locate file info.")
	}
	if len(encryptedFile) < 65 {
		return invitationPtr, errors.New("Data dont make sense bruh.")
	}
	fileHMAC, err := userlib.HMACEval(ownerFileAccess.MAC, encryptedFile[64:])
	if err != nil {
		return invitationPtr, err
	}
	equal = userlib.HMACEqual(fileHMAC, encryptedFile[:64])
	if !equal {
		return invitationPtr, errors.New("Tampering of data detected! Cannot retrieve file.")
	}
	encryptedData = encryptedFile[64:]
	decrypted = userlib.SymDec(ownerFileAccess.Key, encryptedData)
	err = json.Unmarshal(decrypted, &file)
	if err != nil {
		return invitationPtr, err
	}
	
	// Editing SharedUsersMap
	file.SharedUserMap[userdata.Username] = append(file.SharedUserMap[userdata.Username], recipientUsername)
	
	// Encrypting file struct
	marshaled_file, err := json.Marshal(&file)
	if err != nil {
		return invitationPtr, err
	}
	randomFileIV := userlib.RandomBytes(16)
	encrytedFileData := userlib.SymEnc(ownerFileAccess.Key, randomFileIV, marshaled_file)
	hmacFileData, err := userlib.HMACEval(ownerFileAccess.MAC, encrytedFileData)
	if err != nil {
		return invitationPtr, err
	}
	fileStorage := append(hmacFileData, encrytedFileData...)
	userlib.DatastoreSet(file.ID, fileStorage)

	// Populating FileAccess struct
	fileAccess.FileLoc = ownerFileAccess.FileLoc
	fileAccess.Key = ownerFileAccess.Key
	fileAccess.MAC = ownerFileAccess.MAC

	// Encrypting FileAccess struct
	marshaled_fileAcess, err := json.Marshal(&fileAccess)
	if err != nil {
		return invitationPtr, err
	}
	salt = append(userlib.Hash([]byte(recipientUsername)), userlib.Hash([]byte(userdata.Username))...)
	deterministicK = userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key = deterministicK[:16]
	mac = deterministicK[16:]
	randomFileIV = userlib.RandomBytes(16)
	encryptedFileAccessData := userlib.SymEnc(key, randomFileIV, marshaled_fileAcess)
	hmacFileAccess, err := userlib.HMACEval(mac, encryptedFileAccessData)
	if err != nil {
		return invitationPtr, err
	}
	fileAccessStorage := append(hmacFileAccess, encryptedFileAccessData...)

	// Generating a signed RSA encrypted key
	pubKey, ok := userlib.KeystoreGet(recipientUsername + "0")
	if !ok {
		return invitationPtr, errors.New("Recipient user does not exist! Could not fetch their public key.")
	}
	cipherKey, err := userlib.PKEEnc(pubKey, deterministicK)
	if err != nil {
		return invitationPtr, err
	}
	privSign, err := userlib.DSSign(userdata.DSkey, cipherKey)
	if err != nil {
		return invitationPtr, err
	}

	// Populating Invite struct and sending to DS
	invite.EncryptedFileAccess = fileAccessStorage
	invite.EncryptedKey = cipherKey
	invite.SignedMessage = privSign
	marshaled_invite, err := json.Marshal(&invite)
	if err != nil {
		return invitationPtr, err
	}
	invitationPtr, err = uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(recipientUsername))[:8]...))
	if err != nil {
		return invitationPtr, err
	}
	userlib.DatastoreSet(invitationPtr, marshaled_invite)
	return invitationPtr, nil
}


func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	var invite Invite
	var fileAccess FileAccess
	var fileptr File

	if senderUsername == userdata.Username {
		return errors.New("Cannot accept invite to one's own file.")
	}

	// Checking if caller already has filename in their namespace
	ID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return err
	}
	_, ok := userlib.DatastoreGet(ID)
	if ok {
		return errors.New("File with given filename already exists in personal namespace.")
	}

	// Getting invitation and de-crypting it
	marshaled_invite, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New("Could not fetch invite from DataStore. Invitation no longer valid.")
	}
	err = json.Unmarshal(marshaled_invite, &invite)
	if err != nil {
		return err
	}
		//	Verifying sign
	DSVerifyKey, ok := userlib.KeystoreGet(senderUsername + "1")
	if !ok {
		return errors.New("Could not fetch verify key from KeyStore")
	}
	err = userlib.DSVerify(DSVerifyKey, invite.EncryptedKey, invite.SignedMessage)
	if err != nil {
		return err
	}
		// De-crypting key
	plaintxt, err := userlib.PKEDec(userdata.Pkey, invite.EncryptedKey)
	if err != nil {
		return err
	}
	key := plaintxt[:16]
	mac := plaintxt[16:]

		// Decrypting FileAccess struct
	if len(invite.EncryptedFileAccess) < 65 {
		return errors.New("Data dont make sense bruh.")
	}
	hmac, err := userlib.HMACEval(mac, invite.EncryptedFileAccess[64:])
	if err != nil {
		return err
	}
	equal := userlib.HMACEqual(hmac, invite.EncryptedFileAccess[:64])
	if !equal {
		return errors.New("Tampering of data detected! Cannot retrieve file access.")
	}
	encryptedData := invite.EncryptedFileAccess[64:]
	decrypted := userlib.SymDec(key, encryptedData)
	err = json.Unmarshal(decrypted, &fileAccess)
	if err != nil {
		return err
	}
	
	// Decrypting File struct
	encryptedFile, ok := userlib.DatastoreGet(fileAccess.FileLoc)
	if !ok {
		return errors.New("Could not locate file info.")
	}
	if len(encryptedFile) < 65 {
		return errors.New("Data dont make sense bruh.")
	}
	fileHMAC, err := userlib.HMACEval(fileAccess.MAC, encryptedFile[64:])
	if err != nil {
		return err
	}
	equal = userlib.HMACEqual(fileHMAC, encryptedFile[:64])
	if !equal {
		return errors.New("Tampering of data detected! Cannot retrieve file.")
	}
	encryptedData = encryptedFile[64:]
	decrypted = userlib.SymDec(fileAccess.Key, encryptedData)
	err = json.Unmarshal(decrypted, &fileptr)
	if err != nil {
		return err
	}

	// Editing maps and re-encrypting using owner's public RSA
	pubKey, ok := userlib.KeystoreGet(fileptr.Owner + "0")
	if !ok {
		return errors.New("Could not fetch public key from KeyStore")
	}
	encryptedFileName, err := userlib.PKEEnc(pubKey, []byte(filename))
	if err != nil {
		return err
	}
	fileptr.UserFileMap[userdata.Username] = encryptedFileName
	
	// Encrypting File struct
	marshaled_file, err := json.Marshal(&fileptr)
	if err != nil {
		return err
	}
	randomFileIV := userlib.RandomBytes(16)
	encrytedFileData := userlib.SymEnc(fileAccess.Key, randomFileIV, marshaled_file)
	hmacFileData, err := userlib.HMACEval(fileAccess.MAC, encrytedFileData)
	if err != nil {
		return err
	}
	fileStorage := append(hmacFileData, encrytedFileData...)
	userlib.DatastoreSet(fileptr.ID, fileStorage)

	// Encrypting FileAccess struct and re-storing on DataStore
	fileAccessID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return err
	}
	marshaled_fileAcess, err := json.Marshal(&fileAccess)
	if err != nil {
		return err
	}
	salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key = deterministicK[:16]
	mac = deterministicK[16:]
	randomFileIV = userlib.RandomBytes(16)
	encryptedFileAccessData := userlib.SymEnc(key, randomFileIV, marshaled_fileAcess)
	hmacFileAccess, err := userlib.HMACEval(mac, encryptedFileAccessData)
	if err != nil {
		return err
	}
	fileAccessStorage := append(hmacFileAccess, encryptedFileAccessData...)
	userlib.DatastoreSet(fileAccessID, fileAccessStorage)
	
	//Deleting invite from DataStore
	userlib.DatastoreDelete(invitationPtr)

	return nil
}


func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	var fileAccessptr FileAccess
	var fileptr File

	// Error check: If user revoking access to themself
	if userdata.Username == recipientUsername {
		return errors.New("Cannot revoke access on oneself.")
	}
	
	fileAccessID, err := uuid.FromBytes(append(userlib.Hash([]byte(userdata.Username))[:8], userlib.Hash([]byte(filename))[:8]...))
	if err != nil {
		return err
	}
	encryptedFileAccess, ok := userlib.DatastoreGet(fileAccessID)
	if !ok {
		return errors.New("Could not locate encrypted file struct.")
	}

	// Decrypting fileAccess struct	
	salt := append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key := deterministicK[:16]
	mac := deterministicK[16:]
	if len(encryptedFileAccess) < 65 {
		return errors.New("Data dont make sense bruh.")
	}
	hmac, err := userlib.HMACEval(mac, encryptedFileAccess[64:])
	if err != nil {
		return err
	}
	equal := userlib.HMACEqual(hmac, encryptedFileAccess[:64])
	if !equal {
		return errors.New("Tampering of data detected! Cannot retrieve file access.")
	}
	encryptedData := encryptedFileAccess[64:]
	decrypted := userlib.SymDec(key, encryptedData)
	err = json.Unmarshal(decrypted, &fileAccessptr)
	if err != nil {
		return err
	}

	// Decrypting file struct
	encryptedFile, ok := userlib.DatastoreGet(fileAccessptr.FileLoc)
	if !ok {
		return errors.New("Could not locate file info.")
	}
	if len(encryptedFileAccess) < 65 {
		return errors.New("Data dont make sense bruh.")
	}
	fileHMAC, err := userlib.HMACEval(fileAccessptr.MAC, encryptedFile[64:])
	if err != nil {
		return err
	}
	equal = userlib.HMACEqual(fileHMAC, encryptedFile[:64])
	if !equal {
		return errors.New("Tampering of data detected! Cannot retrieve file.")
	}
	encryptedData = encryptedFile[64:]
	decrypted = userlib.SymDec(fileAccessptr.Key, encryptedData)
	err = json.Unmarshal(decrypted, &fileptr)
	if err != nil {
		return err
	}

	// Deleting list of revoked users
	revokeList := RevokedUsers(fileptr.SharedUserMap, recipientUsername)
	for _, user := range revokeList {
		// Check to see if user is invited or accepted
		root_user := TraceValueToKey(fileptr.SharedUserMap, user)
		invitationPtr, err := uuid.FromBytes(append(userlib.Hash([]byte(root_user))[:8], userlib.Hash([]byte(user))[:8]...))
		if err != nil {
			return err
		}
		_, ok := userlib.DatastoreGet(invitationPtr)
		if ok {
			// Delete user invite from DS, delete user from sharedMap
			userlib.DatastoreDelete(invitationPtr)
			fileptr.SharedUserMap[root_user] = DeleteUserFromList(user, fileptr.SharedUserMap[root_user])
		} else {
			// Delete user fileAccess from DS, delete user from userMap, delete user from shared map
			encryptedFilename := fileptr.UserFileMap[user]
			filenamed, err := userlib.PKEDec(userdata.Pkey, encryptedFilename)
			if err != nil {
				return err
			}
			ID, err := uuid.FromBytes(append(userlib.Hash([]byte(user))[:8], userlib.Hash([]byte(filenamed))[:8]...))
			if err != nil {
				return err
			}
			userlib.DatastoreDelete(ID)
			delete(fileptr.UserFileMap, user)

			// TODO: DELETE FROM SHARED MAP
			// shared map becomes a problem when user has shared with other users.
		}
	}

	// Generating new keys, UUID for file struct, Keep track of oldKey for decrypting fileParts
	FileLoc := uuid.New()
	Key := userlib.RandomBytes(16)
	MAC := userlib.RandomBytes(16)
	OldKey := fileAccessptr.Key

	// Updating owner fileAccess struct and file struct
	fileptr.ID = FileLoc
	fileAccessptr.FileLoc = FileLoc
	fileAccessptr.Key = Key
	fileAccessptr.MAC = MAC
	
	// Updating info for all non-revoked users
	for user, value := range fileptr.UserFileMap {
		var userFileAccess FileAccess

		// Decrypt user filename
		filenamed, err := userlib.PKEDec(userdata.Pkey, value)
		if err != nil {
			return err
		}

		// Decrypt user fileAccess struct
		ID, err := uuid.FromBytes(append(userlib.Hash([]byte(user))[:8], userlib.Hash([]byte(filenamed))[:8]...))
		if err != nil {
			return err
		}
		encryptedFileAccess, ok := userlib.DatastoreGet(ID)
		if !ok {
			return errors.New("Could not find encrypted user file access.")
		}
		salt := append(userlib.Hash([]byte(user)), userlib.Hash([]byte(filenamed))...)
		deterministicK := userlib.Argon2Key(userlib.Hash([]byte(filenamed)), salt, 32)
		key := deterministicK[:16]
		mac := deterministicK[16:]
		if len(encryptedFileAccess) < 65 {
			return errors.New("Data dont make sense bruh.")
		}
		hmac, err := userlib.HMACEval(mac, encryptedFileAccess[64:])
		if err != nil {
			return err
		}
		equal := userlib.HMACEqual(hmac, encryptedFileAccess[:64])
		if !equal {
			return errors.New("Tampering of data detected! Cannot retrieve file access.")
		}
		encryptedData := encryptedFileAccess[64:]
		decrypted := userlib.SymDec(key, encryptedData)
		err = json.Unmarshal(decrypted, &userFileAccess)
		if err != nil {
			return err
		}

		// Change user fileAccess struct values
		userFileAccess.FileLoc = FileLoc
		userFileAccess.Key = Key
		userFileAccess.MAC = MAC

		// Encrypt user FileAccess struct
		marshaled_fileAcess, err := json.Marshal(&userFileAccess)
		if err != nil {
			return err
		}
		randomFileIV := userlib.RandomBytes(16)
		encryptedFileAccessData := userlib.SymEnc(key, randomFileIV, marshaled_fileAcess)
		hmacFileAccess, err := userlib.HMACEval(mac, encryptedFileAccessData)
		if err != nil {
			return err
		}
		fileAccessStorage := append(hmacFileAccess, encryptedFileAccessData...)
		userlib.DatastoreSet(ID, fileAccessStorage)
	}

	// Storing fileParts in new locations
	partOldKey := OldKey
	var partOldMAC []byte
	currentUUID := fileptr.Next

	// Update fileptr.Next to new value
	fileptr.Next = uuid.New()
	newUUID := fileptr.Next
	partNewKey := Key
	var partNewMAC []byte

	for partCounter := 0; partCounter < fileptr.Parts; partCounter++ {
		var filePart FilePart
		derivedHashKDF, err := userlib.HashKDF(partOldKey, []byte("Create chain of file nodes encryption"))
		if err != nil {
			return err
		}
		partOldKey = derivedHashKDF[:16]
		partOldMAC = derivedHashKDF[16:32]
		// Decrypt filePart
		encryptedNode, ok := userlib.DatastoreGet(currentUUID)
		if !ok {
			return errors.New("Could not locate file info.")
		}
		if len(encryptedNode) < 65 {
			return errors.New("Data dont make sense bruh.")
		}
		nodeHMAC, err := userlib.HMACEval(partOldMAC, encryptedNode[64:])
		if err != nil {
			return err
		}
		equal = userlib.HMACEqual(nodeHMAC, encryptedNode[:64])
		if !equal {
			return errors.New(fmt.Sprintf("failed at node %d", partCounter))
		}
		encryptedData = encryptedNode[64:]
		decrypted = userlib.SymDec(partOldKey, encryptedData)
		err = json.Unmarshal(decrypted, &filePart)
		if err != nil {
			return err
		}
		// Update filePart next, last, lastNode pointer
		tempUUID := filePart.Next
		filePart.Next = uuid.New()
		// Re-encrypt filePart with new keys
		marshaled_filePart, err := json.Marshal(&filePart)
		if err != nil {
			return err
		}
		derivedHashKDF, err = userlib.HashKDF(partNewKey, []byte("Create chain of file nodes encryption"))
		if err != nil {
			return err
		}
		partNewKey = derivedHashKDF[:16]
		partNewMAC = derivedHashKDF[16:32]
		randomFileIV := userlib.RandomBytes(16)
		encrytedFilePartData := userlib.SymEnc(partNewKey, randomFileIV, marshaled_filePart)
		hmacFilePartData, err := userlib.HMACEval(partNewMAC, encrytedFilePartData)
		if err != nil {
			return err
		}
		filePartStorage := append(hmacFilePartData, encrytedFilePartData...)
		userlib.DatastoreSet(newUUID, filePartStorage)
		// Update looping values
		currentUUID = tempUUID
		newUUID = filePart.Next
		/* TODO: HOW TO UPDATE FILE STRUCT LAST AND LASTNODE POINTERS?
		if partCounter == fileptr.Parts - 2 {
			fileptr.LastFileNode = 
		}*/
	}

	// Re-Encrypting file struct
	marshaled_file, err := json.Marshal(fileptr)
	if err != nil {
		return err
	}
	randomFileIV := userlib.RandomBytes(16)
	encrytedFileData := userlib.SymEnc(fileAccessptr.Key, randomFileIV, marshaled_file)
	hmacFileData, err := userlib.HMACEval(fileAccessptr.MAC, encrytedFileData)
	if err != nil {
		return err
	}
	fileStorage := append(hmacFileData, encrytedFileData...)
	userlib.DatastoreSet(fileptr.ID, fileStorage)

	// Re-Encrypting ownerFileAccess struct
	marshaled_fileAcess, err := json.Marshal(&fileAccessptr)
	if err != nil {
		return err
	}
	salt = append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...)
	deterministicK = userlib.Argon2Key(userlib.Hash([]byte(filename)), salt, 32)
	key = deterministicK[:16]
	mac = deterministicK[16:]
	randomFileIV = userlib.RandomBytes(16)
	encryptedFileAccessData := userlib.SymEnc(key, randomFileIV, marshaled_fileAcess)
	hmacFileAccess, err := userlib.HMACEval(mac, encryptedFileAccessData)
	if err != nil {
		return err
	}
	fileAccessStorage := append(hmacFileAccess, encryptedFileAccessData...)
	userlib.DatastoreSet(fileAccessID, fileAccessStorage)

	return nil
}

// 
// HELPER FUNCS
//

func RevokedUsers(sharedUserMap map[string][]string, mainUser string) []string {
	revokedUsersList := make([]string, 0)
	revokedUsersList = append(revokedUsersList, mainUser)
	for _, user := range sharedUserMap[mainUser] {
		revokedUsersList = append(revokedUsersList, RevokedUsers(sharedUserMap, user)...)
	}
	return revokedUsersList
}

func TraceValueToKey(sharedMap map[string][]string, value string) string {
	for user, list := range sharedMap {
		if InList(value, list) {
			return user
		}
	}
	return ""
}

func InList(value string, lst []string) bool {
	for _, val := range lst {
		if value == val {
			return true
		}
	}
	return false
}

func DeleteUserFromList(user string, list []string) []string {
	for index, value := range list {
		if value == user {
			return append(list[:index], list[index + 1:]...)
		}
	}
	return list
}
