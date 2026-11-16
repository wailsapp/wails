package codesignature

import (
	"bytes"
	"crypto/x509"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Go's Darwin system verifier uses an SSL policy before checking the requested
// EKU. A timestamp-only certificate must instead use Apple's timestamp policy.
// Explicit custom roots retain Go's existing verifier and caller-selected trust.
func verifyTimestampChain(certificate *x509.Certificate, certificates []*x509.Certificate, roots *x509.CertPool, at time.Time) error {
	if roots != nil {
		return verifyTimestampGoChain(certificate, certificates, roots, at)
	}
	api, err := loadTimestampSecurity()
	if err != nil {
		return err
	}
	chain := []*x509.Certificate{certificate}
	for _, candidate := range certificates {
		if !bytes.Equal(candidate.Raw, certificate.Raw) {
			chain = append(chain, candidate)
		}
	}
	refs := make([]uintptr, 0, len(chain))
	for _, cert := range chain {
		if len(cert.Raw) == 0 {
			return errors.New("empty timestamp certificate")
		}
		data := api.dataCreate(0, unsafe.Pointer(&cert.Raw[0]), len(cert.Raw))
		runtime.KeepAlive(cert.Raw)
		if data == 0 {
			return errors.New("could not allocate timestamp certificate data")
		}
		ref := api.certificateCreate(0, data)
		api.release(data)
		if ref == 0 {
			return errors.New("Security.framework rejected timestamp certificate encoding")
		}
		defer api.release(ref)
		refs = append(refs, ref)
	}
	array := api.arrayCreate(0, unsafe.Pointer(&refs[0]), len(refs), api.arrayCallbacks)
	runtime.KeepAlive(refs)
	if array == 0 {
		return errors.New("could not create timestamp certificate chain")
	}
	defer api.release(array)
	policy := api.policyCreate(api.timestampPolicy, 0)
	if policy == 0 {
		return errors.New("could not create platform timestamp trust policy")
	}
	defer api.release(policy)
	var trust uintptr
	if status := api.trustCreate(array, policy, &trust); status != 0 {
		return fmt.Errorf("create timestamp trust: OSStatus %d", status)
	}
	if trust == 0 {
		return errors.New("platform returned an empty timestamp trust")
	}
	defer api.release(trust)
	date := api.dateCreate(0, float64(at.Unix()-978307200)+float64(at.Nanosecond())/1e9)
	if date == 0 {
		return errors.New("could not create timestamp verification date")
	}
	defer api.release(date)
	if status := api.setDate(trust, date); status != 0 {
		return fmt.Errorf("set timestamp verification date: OSStatus %d", status)
	}
	// All submitted intermediates are supplied. Avoid starting unbounded network
	// work inside a synchronous platform call; cached/system trust still applies.
	if status := api.setNetworkFetch(trust, false); status != 0 {
		return fmt.Errorf("set timestamp trust network policy: OSStatus %d", status)
	}
	var failure uintptr
	valid := api.evaluate(trust, &failure)
	if failure != 0 {
		defer api.release(failure)
	}
	if !valid {
		code := int64(0)
		if failure != 0 {
			code = api.errorCode(failure)
		}
		return fmt.Errorf("platform timestamp trust rejected certificate chain (code %d)", code)
	}
	return nil
}

type timestampSecurity struct {
	release                         func(uintptr)
	dataCreate                      func(uintptr, unsafe.Pointer, int) uintptr
	arrayCreate                     func(uintptr, unsafe.Pointer, int, uintptr) uintptr
	dateCreate                      func(uintptr, float64) uintptr
	errorCode                       func(uintptr) int64
	certificateCreate               func(uintptr, uintptr) uintptr
	policyCreate                    func(uintptr, uintptr) uintptr
	trustCreate                     func(uintptr, uintptr, *uintptr) int32
	setDate                         func(uintptr, uintptr) int32
	setNetworkFetch                 func(uintptr, bool) int32
	evaluate                        func(uintptr, *uintptr) bool
	arrayCallbacks, timestampPolicy uintptr
}

var loadTimestampSecurity = sync.OnceValues(func() (*timestampSecurity, error) {
	core, err := purego.Dlopen("/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation", purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		return nil, err
	}
	security, err := purego.Dlopen("/System/Library/Frameworks/Security.framework/Security", purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		_ = purego.Dlclose(core)
		return nil, err
	}
	ok := false
	defer func() {
		if !ok {
			_ = purego.Dlclose(security)
			_ = purego.Dlclose(core)
		}
	}()
	api := &timestampSecurity{}
	for _, binding := range []struct {
		handle uintptr
		name   string
		target any
	}{
		{core, "CFRelease", &api.release},
		{core, "CFDataCreate", &api.dataCreate},
		{core, "CFArrayCreate", &api.arrayCreate},
		{core, "CFDateCreate", &api.dateCreate},
		{core, "CFErrorGetCode", &api.errorCode},
		{security, "SecCertificateCreateWithData", &api.certificateCreate},
		{security, "SecPolicyCreateWithProperties", &api.policyCreate},
		{security, "SecTrustCreateWithCertificates", &api.trustCreate},
		{security, "SecTrustSetVerifyDate", &api.setDate},
		{security, "SecTrustSetNetworkFetchAllowed", &api.setNetworkFetch},
		{security, "SecTrustEvaluateWithError", &api.evaluate},
	} {
		address, err := purego.Dlsym(binding.handle, binding.name)
		if err != nil {
			return nil, err
		}
		purego.RegisterFunc(binding.target, address)
	}
	api.arrayCallbacks, err = purego.Dlsym(core, "kCFTypeArrayCallBacks")
	if err != nil {
		return nil, err
	}
	lookupAddress, err := purego.Dlsym(purego.RTLD_DEFAULT, "dlsym")
	if err != nil {
		return nil, err
	}
	// Keep the exported variable's address as a pointer while loading its value.
	var lookupVariable func(uintptr, string) unsafe.Pointer
	purego.RegisterFunc(&lookupVariable, lookupAddress)
	policyAddress := lookupVariable(security, "kSecPolicyAppleTimeStamping")
	if policyAddress == nil {
		return nil, errors.New("platform timestamp policy symbol is unavailable")
	}
	api.timestampPolicy = *(*uintptr)(policyAddress)
	if api.timestampPolicy == 0 {
		return nil, errors.New("platform timestamp policy identifier is unavailable")
	}
	ok = true // Function pointers remain valid for the process lifetime.
	return api, nil
})
