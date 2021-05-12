package domain

// HintTelnetFailed invoked when expect was failed
const HintTelnetFailed = `
[Hint]
- Are your username and password correct?
  Some environments may use local and LDAP accounts.
- Does the set host name match the actual host name?
  Some devices require the host name for expectation.
- Are the exec-platform you set correct?
  Some devices also needs --ha-mode option.
`
