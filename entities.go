package corezoid

type ContentType int

const (
	Unknown ContentType = iota
	Json
)

type Op map[string]interface{}

func (o Op) IsOK() bool {
	if v, ok := o["proc"].(string); ok {
		return v == "ok"
	} else {
		return false
	}
}

type Ops struct {
	List []Op `json:"ops"`
}

func (o *Ops) Add(op Op) {
	if op == nil {
		return
	}
	o.List = append(o.List, op)
}

type OpsResult struct {
	StatusCode  int
	RequestProc string `json:"request_proc"`
	List        []Op   `json:"ops"`
}

func (o *OpsResult) IsSuccessCode() bool {
	return o.StatusCode >= 200 || o.StatusCode < 300
}

func (o *OpsResult) IsRequestProcOK() bool {
	return o.RequestProc == "ok"
}

func (o *OpsResult) IsRequestOK() bool {
	return o.IsRequestProcOK() && o.IsSuccessCode()
}

func (o *OpsResult) IsOpsOK() bool {
	for _, x := range o.List {
		if !x.IsOK() {
			return false
		}
	}

	return true
}

type Task struct {
	Type   string                 `json:"type"`
	ConvID int                    `json:"conv_id"`
	Obj    string                 `json:"obj"`
	Data   map[string]interface{} `json:"data"`
	Ref    string                 `json:"ref,omitempty"`
	ID     string                 `json:"id,omitempty"`
}

func NewTask(convID int, ref string) Task {
	return Task{
		Type:   "create",
		Obj:    "task",
		ConvID: convID,
		Ref:    ref,
		Data:   make(map[string]interface{}),
	}
}

func (t *Task) SetID(ID string) {
	t.ID = ID
}

func (t *Task) SetData(data map[string]interface{}) {
	t.Data = data
}

func (t *Task) Put(key string, value interface{}) {
	t.Data[key] = value
}
