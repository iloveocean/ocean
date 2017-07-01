package mongo

type Condition struct {
	conditions Map //map to store query conditions
}

func NewCondition() *Condition {
	value := &Condition{
		conditions: Map{},
	}
	return value
}

func isLogicQuery(opt Operator) bool {
	if opt == Or || opt == And || opt == Nor || opt == ElemMatch {
		return true
	}
	return false
}

func (c *Condition) addCond(key string, value interface{}, opt Operator) *Condition {
	if isLogicQuery(opt) {
		newOne := M(key, value)
		c.conditions.MSlice(string(opt), newOne)
	} else {
		c.conditions.M(key, value, opt)
	}
	return c
}

func (c *Condition) mergeCond(newConds *Condition) *Condition {
	for k, v := range newConds.conditions {
		c.conditions[k] = v
	}
	return c
}
