- **Secure Remote Execution**: Uses gRPC authentication to securely manage remote Docker instances.
- **Automation-Friendly**: Enables running scripts and commands on containers programmatically.
- **File Transfers**: Easily push and pull files between local machines and containers.
- **Efficient CLI Interface**: Built with Bubble Tea for a user-friendly terminal experience.
- **Lightweight and Extensible**: Designed for performance with minimal dependencies.
- **Cross-Platform Support**: Works on Windows, macOS, and Linux.
- gRPC-based communication for efficient and secure remote management, fast file transfers, and real-time command execution.
- **Secure authentication** using Kerberos, OAuth 2.0, or SAML
- **File transfer support** with real-time progress updates
- **Command execution** with streaming output

- Optimize file transfer protocols to reduce latency.
- Implement real-time logs streaming from remote containers.

kerberos:
### 🔹 **Client Side:**

1. Client authenticates with the **Kerberos Authentication Server (AS)** and receives a **Ticket-Granting Ticket (TGT)**.
2. Client uses the TGT to request a **Service Ticket (ST)** from the **Ticket Granting Server (TGS)**.
3. Client sends the ST in the **gRPC request** as metadata to the server.

### 🔸 **Server Side:**

1. **First Request**:
    
    - Server validates the ST locally using the Kerberos shared key (without contacting the KDC).
    - If the ST is valid, process the request.
    - Cache the ST using its expiration time as the TTL.
2. **Subsequent Requests**:
    
    - If the ST is cached and not expired → No need to validate again.
    - If the ST is expired → Request a new ST from the client.
3. **Ticket Renewal (Optional)**:
    - Background worker can pre-fetch renewable tickets to reduce latency when the ticket expires.

oauth:
### 🔹 **Client Side:**

1. Client authenticates with the **OAuth Authorization Server** and receives an **Access Token** (JWT or opaque).
2. Client includes the token in the **gRPC request** as metadata.

### 🔸 **Server Side:**

1. **Fetching Signing Keys**:
    
    - On startup, the server fetches the signing keys (via the JWKS endpoint).
    - Cache the keys with a TTL based on the provider’s metadata (or a fixed interval).
2. **First Request**:
    
    - If token is a JWT → Validate locally using the cached signing key.
3. **Subsequent Requests**:  
    ✅ **JWT:**
    
    - If token is a JWT and cached → No need to re-validate.
    - If signature or expiration fails → Reject the request.
    
4. **Key Rotation**:
    - Background worker periodically refreshes signing keys.
    - If the token signature fails due to rotation, re-fetch the keys.

saml:
### 🔹 **Client Side:**

1. Client logs into the **Identity Provider (IdP)** and receives a **SAML Assertion**.
2. Client sends the assertion in the **gRPC request** as metadata.

### 🔸 **Server Side:**

1. **Fetching Signing Keys**:
    - On startup, the server fetches the IdP’s signing key (via metadata endpoint).
    - Cache the key and rotate periodically based on IdP metadata.
2. **First Request**:
    
    - Validate the SAML assertion’s signature using the cached key.
    - Verify the `NotOnOrAfter` attribute (expiration).
    - If valid → Convert the assertion into a session token and cache it.
3. **Subsequent Requests**:  
    ✅ If session token is cached and valid → Skip SAML validation.  
    ✅ If assertion is cached and valid → No need to contact IdP.  
    ✅ If session expires → Refresh by validating a new assertion.
    
4. **Session Renewal**:
    - Background worker can refresh long-lived sessions before expiration to avoid latency.
    - If the IdP’s key rotates → Refresh key cache automatically.

- The server periodically refreshes the JWKS key set using a background worker.
- If signature verification fails → Attempt to fetch a fresh key set.
### ✅ **Why It’s Efficient:**

- JWT = Local validation → No need to contact the provider per request.
- Key rotation is handled in the background → No client-side impact.
## 🔑 **How Kerberos Works**

1. **Client Authentication Phase**
    
    - The client contacts the **Kerberos Authentication Server (AS)** and requests a **Ticket-Granting Ticket (TGT)** by providing credentials (like username and password).
    - The AS encrypts the TGT using the **Kerberos secret key** shared with the **Ticket Granting Server (TGS)** and sends it back to the client.
2. **Service Ticket Request Phase**
    
    - The client sends the TGT to the **TGS** and requests a **Service Ticket (ST)** for the specific service it wants to access (in this case, your server).
    - The TGS issues the ST, encrypted using the **service’s shared key** (which is known to both the TGS and the server).
3. **Client Access Phase**
    
    - The client sends the ST to the server in the gRPC request.
    - The server uses its own **Kerberos shared key** (pre-configured during Kerberos setup) to decrypt and validate the ST.
    - If the ticket is valid → The request is authenticated and processed.

## 🔑 **How SAML Works**

**(Used for Single Sign-On (SSO) in enterprise applications)**

### 1. **Client Authentication Phase**

- The client accesses the service and is redirected to the **Identity Provider (IdP)** (e.g., Okta, Microsoft AD).
- The client logs into the IdP using credentials (e.g., username/password, MFA).
- The IdP verifies the credentials and generates a **SAML Assertion** (an XML document):
    - Signed using the IdP’s private key.
    - Contains user identity, roles, and expiration (`NotOnOrAfter`).
    - Optionally encrypted for added security.

---

### 2. **Assertion Delivery Phase**

- The IdP sends the **SAML Assertion** to the client via a browser redirect (or POST).
- The client forwards the assertion to the server as part of the gRPC request.

---

### 3. **Assertion Validation Phase**

- The server validates the SAML Assertion:  
    ✅ Signature → Validated using the IdP’s public key (retrieved from metadata endpoint).  
    ✅ Expiration → Checked against `NotOnOrAfter` attribute.  
    ✅ User Identity → Extracted and verified against the system’s expected user state.
    
- If valid → The request is authenticated and processed.
    
- The server caches the decoded assertion using the `NotOnOrAfter` time as the TTL.
    

---

### 5. **Subsequent Requests**

- If the session token is cached → Skip SAML validation.
- If the session expires → Request a new SAML assertion from the client.
- If the IdP’s public key rotates → Refresh keys from the metadata endpoint.

---

### 6. **Key Rotation and Refresh**

- The server periodically refreshes the IdP’s public key using a background worker.
- If signature verification fails → Attempt to refresh the key from the IdP.
- If the key endpoint fails due to rate limiting → Implement exponential backoff.

---

## 🔑 **How OAuth 2.0 Works**

**(Used for API-based authentication via Identity Providers like Google, Okta, etc.)**

### 1. **Client Authentication Phase**

- The client sends a request to the **OAuth Authorization Server** (e.g., Google, Okta) to initiate the authentication flow.
- The client presents credentials (such as a client ID, client secret, or user password) depending on the grant type:
    - **Authorization Code** – User login with redirection.
    - **Client Credentials** – Machine-to-machine authentication.
    - **Implicit** – Direct client-side token retrieval.
    - **Refresh Token** – Long-term session maintenance.

---

### 2. **Access Token Request Phase**

- If the credentials are valid, the Authorization Server issues an **Access Token**:  
    ✅ If the token is a **JWT** → It is signed using the provider’s private key.  

---

### 3. **Client Access Phase**

- The client includes the **Access Token** in the gRPC request as metadata (`Authorization: Bearer <token>`).
    
- The server handles the token differently based on its type:  
    ✅ **JWT:**
    
    - The server verifies the token’s signature using the provider’s public key (retrieved from the JWKS endpoint).
    - If the signature is valid and the token is not expired → The request is authenticated and processed.
    

### 4. **Subsequent Requests**

- **JWT:**  
    ✅ If the token is cached → Skip verification until expiration.  
    ✅ If signature verification fails → Refresh keys from JWKS endpoint.
    
---

### 5. **Key Rotation and Refresh**

- The server periodically refreshes the JWKS key set using a background worker.
- If signature verification fails → Attempt to fetch a fresh key set.
- If the introspection endpoint fails due to rate limiting → Implement exponential backoff.

---

### ✅ **Why It’s Efficient:**

- JWT = Local validation → No need to contact the provider per request.
- Key rotation is handled in the background → No client-side impact.