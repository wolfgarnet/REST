package REST

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"web/core/response"
	"github.com/wolfgarnet/stopwatch"
	"github.com/wolfgarnet/logging"
)

var logger logging.Logger

type RestHandler struct {
	System System
	root Node
}

func NewRestHandler(system System, root Node) {
	return &RestHandler{system, root}
}

func doRedirect(a Autonomous, rest string, endsWithSlash bool) bool {
	if a != nil {
		return false
	}

	if rest == "" && !endsWithSlash {
		return true
	}

	return false
}

func (rh *RestHandler) resolve(parts []string) (Node, string, Autonomous) {

	var current Node = rh.root
	var next Node = nil
	var err error

	//var token = parts[0]
	var rest = parts[1]
	fmt.Println("LEN", len(parts))
	for len(parts[1]) > 0 {
		parts = strings.SplitAfterN(rest, "/", 2)
		log.Printf("PARTS: %v", parts)
		log.Printf("Current: %v", current)
		log.Printf("REST: %v", rest)

		token := strings.TrimRight(parts[0], "/")
		fmt.Println("TOKEN:", token)

		// Reset
		next = nil

		// Autonomous
		a, ok := current.(Autonomous)
		if ok {
			log.Printf("Hey, this was autonomous")
			return current, "", a
		}

		parent, ok := current.(Parent)
		if ok {
			logger.Debug("Current is a parent")
			next, err = parent.GetChild(token)
		}

		// If nil , find action
		actions := rh.System.GetExtensions("ProcessAction")
		logger.Debug("I HAVE ACTIONS: {} for {}", actions, token)
		for _, a := range actions {
			logger.Debug("ACTION IS {}", a)
			action, ok := a.(Action)
			if ok && action.GetName() == token {
				logger.Debug("Found action, %v", token)
				if action.IsApplicable(current) {
					logger.Debug("Action is applicable, %v", action)
				}
			}
		}

		// If still nil, bail
		if next == nil || err != nil {
			log.Printf("DOEST, %v --- %v", rest, err)
			break
		}

		// If not nil, verify authorization

		current = next

		// No more parts to parse
		if len(parts) < 2 {
			logger.Debug("HER: %v --- %v", current, parts)
			if current != nil {
				rest = ""
			}
			break;
		}


		//token = parts[0]
		rest = parts[1]
	}

	fmt.Println("The rest", rest)

	return current, rest, nil
}

func (rh *RestHandler) ServeHTTP(rsp http.ResponseWriter, request *http.Request) {
	sw := stopwatch.NewStopWatch()
	sw.Start("rest")

	context := newContext(request, rh.System)

	println("------------------------------")
	println("Processing:", request.URL.String())
	println(context.Session)


	runner := rh.getRunner(context)
	if runner != nil {
		r, err := runner.Run(context)
		if err != nil {
			r = MakeErrorResponse(err.Error())
		}

		logger.Debug("Runner: %v", r)
		rsp.Header().Add("Content-type", r.GetContentType())

		if r.GetStatus() == 200 {
			r.WriteBody(rsp, request)
		} else {
			rsp.WriteHeader(r.GetStatus())
			r.WriteBody(rsp, request)
		}

	} else {
		response.GenerateErrorResponse(rsp, "Error, no such response", 500)
	}

	sw.Stop("rest")

	sw.Print()
}

func (rh *RestHandler) getRunner(context *Context) Runner {

	//fmt.Fprintf(response, "<h1>%s</h1><div>%s</div>", "YEA", "SNADE")
	parts := strings.SplitAfterN(context.Request.URL.Path, "/", 2)

	node, rest, autonomous := rh.resolve(parts)
	//log.Printf("REDIRECTING: %v, %v, %v", rest, autonomous, context.Request.URL.String()[len(context.Request.URL.String())-1:] == "/")
	if doRedirect(autonomous, rest, context.Request.URL.String()[len(context.Request.URL.String())-1:] == "/") {
		return &Redirect{context.Request.URL.String() + "/"}
	}

	if autonomous != nil {
		return &AutonomousRunner{autonomous, context}
	}

	logger.Debug("Last node: %v", node)
	
	if len(rest) > 0 {
		if strings.Contains(rest, "/") {
			return &ErrorRunner{404, "Does not exist: " + rest}
		} else {
			runner := rh.System.GetRunner(context.Session, node, rest, context.Request.Method)
			if runner != nil {
				return runner
			}

			return &ErrorRunner{404, "Does not exist: " + rest}
		}
	} else {
		runner := rh.System.GetRunner(context.Session, node, "index", context.Request.Method)
		if runner != nil {
			return runner
		}

		return &ErrorRunner{404, "Does not exist: " + rest}
	}

	return nil
}
