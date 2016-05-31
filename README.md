# Golang ast-tool

Tool to help parsing your own golang source-code from the ast (=abstract syntax tree) into an intermediate representation.
From this intermediate representation, we can easily generate boring and error-phrone boilerplate source-code.

## Example:

### input:
A regular golang struct definition with our own "+event"-annotation. 
This annotation is used to trigger code-generation

    // +event -> aggregate: tour
    type EtappeCreated struct {
	    Year                 int
	    EtappeId             int
	    EtappeDate           time.Time
	    EtappeStartLocation  strin
	    EtappeFinishLocation string
	    EtappeLength         int
	    EtappeKind           int
    }

### result 
etappeCreatedWrapper.go (approx. 60 lines of code):

    func (s *EtappeCreated) Wrap(aggregateName string, aggegateUid string) (*Envelope,error) {
        ....
    }
    
    func IsEtappeCreated(envelope *Envelope) bool {
        ...
    }

    func GetIfIsEtappeCreated(envelop *Envelope) (*EtappeCreated, bool) {
        ...
    }

    func UnWrapEtappeCreated(envelop *Envelope) (*EtappeCreated,error) {
        ...
    }    
    
