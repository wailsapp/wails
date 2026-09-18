package ast

type Visitor interface {
	VisitTaskfile(*Taskfile) Visitor
	VisitTask(*Task) Visitor
	VisitCmd(*Cmd) Visitor
	VisitDep(*Dep) Visitor
	VisitVar(*Var) Visitor
}

func WalkTaskfile(tf *Taskfile, v Visitor) {
	v = v.VisitTaskfile(tf)
	if v == nil {
		return
	}
	for _, inc := range tf.Includes {
		if inc.Resolved != nil {
			WalkTaskfile(inc.Resolved, v)
		}
	}
	for _, vr := range tf.Vars {
		walkVar(vr, v)
	}
	for _, t := range tf.Tasks {
		WalkTask(t, v)
	}
}

func WalkTask(t *Task, v Visitor) {
	v = v.VisitTask(t)
	if v == nil {
		return
	}
	for _, vr := range t.Vars {
		walkVar(vr, v)
	}
	for _, d := range t.Deps {
		walkDep(d, v)
	}
	for _, c := range t.Cmds {
		walkCmd(c, v)
	}
}

func walkVar(vr *Var, v Visitor) {
	if v.VisitVar(vr) != nil {
	}
}

func walkDep(d *Dep, v Visitor) {
	if v.VisitDep(d) != nil {
	}
}

func walkCmd(c *Cmd, v Visitor) {
	v = v.VisitCmd(c)
	if v == nil {
		return
	}
	if c.For != nil {
		for _, vr := range c.For.Vars {
			walkVar(vr, v)
		}
	}
	for _, vr := range c.Vars {
		walkVar(vr, v)
	}
}
