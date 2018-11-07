package beater

import (
  
  "github.com/elastic/beats/libbeat/common"

  "github.com/go-ole/go-ole"
  "github.com/go-ole/go-ole/oleutil"

)

// WmiQuery is execute query with WMI
func WmiQuery(q string, fields []string) (interface {}, error) {

  ole.CoInitialize(0)

  defer ole.CoUninitialize()

  wmiscriptObj, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
  if err != nil {
    return "", err
  }

  wmiqi, err := wmiscriptObj.QueryInterface(ole.IID_IDispatch)
  if err != nil {
    return "", err
  }
  
  defer wmiscriptObj.Release()
  
  svcObj, err := oleutil.CallMethod(wmiqi, "ConnectServer")
  if err != nil {
    return "", nil
  }

  defer wmiqi.Release()
  
  svc := svcObj.ToIDispatch()
  
  defer svcObj.Clear() 


  rsObj, err := oleutil.CallMethod(svc, "ExecQuery", q)
  if err != nil {
    return "", err
  }

  rs := rsObj.ToIDispatch()
  
  defer rsObj.Clear() 

  cntObj, err := oleutil.GetProperty(rs, "Count")
  if err != nil {
    return "", err
  }

  count := int(cntObj.Val)
  
  defer cntObj.Clear() 

  var clsValues interface {} = nil 
  clsValues = []common.MapStr{}

  for i := 0; i<count; i++ {
    rowObj, err := oleutil.CallMethod(rs, "ItemIndex", i) 
    if err != nil {
      //return "", err
      break
    }

    row := rowObj.ToIDispatch() 
    //defer rowObj.Clear()

    var rowValues common.MapStr

    for _, j := range fields {
      wmiObj, err := oleutil.GetProperty(row, j)
      if err != nil {
        return "", err
      }

      var objValue = wmiObj.Value()

      rowValues = common.MapStrUnion(rowValues, common.MapStr {j: objValue})
      //defer wmiObj.Clear()
      wmiObj.Clear()
    }

    clsValues = append(clsValues.([]common.MapStr), rowValues)
    rowValues = nil
    rowObj.Clear()
  }

  // for new adding swpark
  // cntObj.Clear()
  // rsObj.Clear() 
  // svcObj.Clear()
  // wmiqi.Release()
  // wmiscriptObj.Release()

  //ole.CoUninitialize()

  return clsValues, nil
}