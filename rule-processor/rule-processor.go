package main

import (
	"fmt"
	"time"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type MyFact struct {
    IntAttribute       int64
    StringAttribute    string
    BooleanAttribute   bool
    FloatAttribute     float64
    TimeAttribute      time.Time
    WhatToSay          string
    WhatToSayStatus		string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
    return fmt.Sprintf("Let say \"%s\"", sentence)
}

func main() {

	// Add Fact Into DataContext
	myFact := &MyFact{
	    IntAttribute: 123,
	    StringAttribute: "Some string value",
	    BooleanAttribute: true,
	    FloatAttribute: 1.234,
	    TimeAttribute: time.Now(),
	}

	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("MF", myFact)
	if err != nil {
	    panic(err)
	}


	executeDrlRule(dataCtx, myFact)
}

func executeJsonRule() {

}

func executeDrlRule(dataCtx ast.IDataContext, myFact *MyFact) {
	// Creating a KnowledgeLibrary and Adding Rules Into It
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	// lets prepare a rule definition
	drls := `
	rule CheckValues "Check the default values" salience 10 {
	    when 
	        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
	    then
	        MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
	        Retract("CheckValues");
	}

	rule CheckStatus "Check Status of WhatToSay" salience 10 {
	    when 
	        MF.WhatToSay == MF.GetWhatToSay("Hello Grule")
	    then
	        MF.WhatToSayStatus = "WhatToSay Computed";
	        Retract("CheckStatus");
	}

	`
	
	// Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
	bs := pkg.NewBytesResource([]byte(drls))
	err := ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
	    panic(err)
	}


	// Executing Grule Rule Engine
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")
	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
	    panic(err)
	}

	fmt.Println(myFact.WhatToSay, "\n", myFact.WhatToSayStatus)

}