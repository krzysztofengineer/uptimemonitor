# TODO

What's left to do:

- [ ] Log in after setup automatically
- [ ] Change webhook body available variables (just .Url, .ResponseTimeMs,
      .IncidentUrl, .MonitorUrl, .StatusCode, .Body, .Headers)
- [ ] Add all available variables to the sample webhook body
- [ ] Loading placeholders (skeletons)
- [ ] Test check timeout
- [ ] Calculate overall uptime and response time (keep track of number of checks
      historically)
- [ ] Session expiration
- [ ] Save check response
- [ ] Add check response modal (div#modals + append when clicked + destroy when
      clicked)
- [ ] Save incident request and add a way to preview it
- [ ] Move incident to separate page with request and response details
- [ ] Store request url in incident and check (it could be updated in monitor)
- [ ] Add sponsors badges
- [ ] Add documentation
- [ ] Document how to install
- [ ] Document webhook parsing and available variables
- [ ] Add "Test Webhook" button with fake incident
- [ ] Add website command with a server homepage

Optional:

- [ ] Monitor status
- [ ] Change password
- [ ] Add users
- [ ] Remove users
- [ ] Change user passwords
- [ ] Timezones
- [ ] Reset password via cli
- [ ] Use the same form for creation and editing of monitor
