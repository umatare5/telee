package domain

// HintTelnetFailed is printed on any failed session despite the name, including the SSH-only srx.
const HintTelnetFailed = `
[Hint]
- Are your username and password correct?
  Some environments may use local and LDAP accounts.
- Does the set host name match the actual host name?
  Some devices require the host name for expectation.
- Are the exec-platform you set correct?
  Some devices also needs --redundant-mode option.
`
