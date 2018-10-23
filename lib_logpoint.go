package main

import (
	"encoding/json"
	"errors"
	"fmt"
    "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/core/util"
)

type LibraryUserPoint struct {
}

// ポイント情報
type UserPoint struct {
    Tms string `json: "tms"`
	Name string `json: "name"`
    Pass string `json: "pass"`
	Point int64 `json: "point"`
}


// ユーザー情報の初期値を設定
func (lp *LibraryUserPoint) Init(stub shim.ChaincodeStubInterface, function string, args []string)([]byte,error) {
    if len(args) == 1 {
        f := "get_all"
        queryArgs := util.ToChaincodeArgs(f)
        val, err := stub.QueryChaincode(args[0], queryArgs)
        if err != nil {
            return nil, errors.New("Unable to call chaincode " + args[0])
        }

        var states interface{}
        err = json.Unmarshal(val, &states)
        if err != nil {
            return nil, errors.New("Unable to marshal chaincode return value " + string(val))
        }

        for _, stateIf := range states.([]interface{}){
            state := stateIf.([]interface{})
            stateKey := state[0].(string)
            stateVal := state[1].(string)
            if stateKey == "" {
                return nil, errors.New("Unable to PutState: missing statekey [ ]" + stateKey + " , " + string(stateVal) + " ]")
            }
            err = stub.PutState(stateKey, []byte(stateVal))
            if err != nil {
                return nil, errors.New("Unable to PutState [ " +stateKey + " , " + string(stateVal) + " ]")
            }
        }
    } else {

        //var user UserPoint
        //var userBytes []byte

        //var ntime string
        //ntime = args[0]

        // ユーザー情報を生成
        //user = UserPoint{Tms: ntime, Name: "user0", Pass: "pass0", Point: 0}

        // JSON形式に変換
        //userBytes, _ = json.Marshal(user)
        // ワールドステートに追加
        //stub.PutState(user.Tms, userBytes)

        return nil, nil
    }
    return nil, nil
}


// ユーザー情報を登録
func (lp *LibraryUserPoint) Invoke(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
    var userId string
    var userPs string
    var Addbytes []byte

    var ntime string
    ntime = args[0]

    userId = args[1]
    userPs = args[2]

    auser := UserPoint{}

    auser.Tms = ntime
    auser.Name = userId
    auser.Pass = userPs

    // function名でハンドリング
    if function == "pointUp" {
        addPt, _ := strconv.ParseInt(args[3], 10, 64)
        auser.Point = addPt
    } else if function == "addUser" {
        auser.Point = 0
    }

    // JSON形式に変換
    Addbytes, _ = json.Marshal(auser)
    // ワールドステートに追加
    stub.PutState(auser.Tms, Addbytes)

    return nil, nil
}


// ユーザー情報を参照
func (lp *LibraryUserPoint) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
    // function名でハンドリング
    if function == "get_all" {
        if len(args) != 0 {
            fmt.Printf("Incorrect number of arguments passed")
            return nil, errors.New("QUERY: Incrrect number of arguments passed")
        }
        return lp.get_all(stub)
    }
	if function == "refresh" {
		// カウンター情報を取得
		return lp.getUsers(stub, args)
	}
    if function == "get_log" {
        // ログ情報を取得
        return lp.getLog(stub, args)
    }

	return nil, errors.New("Received unknown function")
}


// ユーザー情報の取得
func (lp *LibraryUserPoint) getUsers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userId string
    var userPs string
    var pointSum int64
var uflg int

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
    }

    userId = args[0]
    userPs = args[1]

    pointSum = 0

uflg = 0

    Valbytes, err := lp.get_all(stub)

    var states interface{}
    err = json.Unmarshal(Valbytes, &states)
    if err != nil {
        return nil, errors.New("Unable to marshal chaincode return value " + string(Valbytes))
    }

    for _, stateIf := range states.([]interface{}){
        state := stateIf.([]interface{})
        stateKey := state[0].(string)
        stateValStr := state[1].(string)
        if stateKey == "" {
            return nil, errors.New("Unable to PutState: missing statekey [ " + stateKey + " ]")
        }

        stateValBytes := ([]byte)(stateValStr)

        userValue := UserPoint{}

        err = json.Unmarshal(stateValBytes, &userValue)
        if err != nil {
            return nil, errors.New("Unable to marshal chaincode return value " + (string)(Valbytes))
        }

        if userId == userValue.Name {
            pointSum = pointSum + userValue.Point
uflg = 1
        }
    }

    if uflg == 0 {
        return nil, errors.New("User dose not exist - " + userId)
    }

    user := UserPoint{}
    user.Tms = ""
    user.Name = userId
    user.Pass = userPs
    user.Point = pointSum

    Pointbytes, _ := json.Marshal(user)

    return Pointbytes, nil
}

//  ログ情報の取得
func (lp *LibraryUserPoint) getLog(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var userId string
    //var userPs string

    var userTupples []UserPoint

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
    }

    userId = args[0]
    //userPs = args[1]

    Valbytes, err := lp.get_all(stub)

    var states interface{}
    err = json.Unmarshal(Valbytes, &states)
    if err != nil {
        return nil, errors.New("Unable to marshal chaincode return value " + string(Valbytes))
    }

    for _, stateIf := range states.([]interface{}){
        state := stateIf.([]interface{})
        stateKey := state[0].(string)
        stateValStr := state[1].(string)
        if stateKey == "" {
            return nil, errors.New("Unable to PutState: missing statekey [ " + stateKey + " ]")
        }

        stateValBytes := ([]byte)(stateValStr)

        userValue := UserPoint{}

        err = json.Unmarshal(stateValBytes, &userValue)
        if err != nil {
            return nil, errors.New("Unable to marshal chaincode return value " + (string)(Valbytes))
        }

        if userId == userValue.Name {
            user := UserPoint{}
            user.Tms = userValue.Tms
            user.Point = userValue.Point

            userTupples = append(userTupples, user)
        }
    }

    marshalledTupples, err := json.Marshal(userTupples)
    return []byte(marshalledTupples), nil
}

func (lp *LibraryUserPoint) get_all(stub shim.ChaincodeStubInterface) ([]byte, error) {
    var tupples [][]string

    keysIter, err := stub.RangeQueryState("", "~")
    if err != nil {
        return nil, errors.New("unable to start the iterator")
    }

    defer keysIter.Close()

    for keysIter.HasNext() {
        key, val, iterErr := keysIter.Next()
        if iterErr != nil {
            return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
        }
        tupple := []string{string(key), string(val)}
        tupples = append(tupples, tupple)
    }

    marshalledTupples, err := json.Marshal(tupples)
    return []byte(marshalledTupples), nil
}


// Validating Peerに接続し、チェーンコードを実行
func main() {
	err := shim.Start(new(LibraryUserPoint))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}





