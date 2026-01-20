This is report given by Gemini on how I solved the problem: Error
sending POST requests from frontend clien to the backend (I don't have
enough knowledge in WebDev)

---

## What Went Wrong?

There were actually two separate "walls" your request was hitting:

### 1. The "Pre-flight" Requirement

When you use JavaScript `fetch` to send JSON (`Content-Type:
application/json`) to a different port, the browser doesn't just send
the `POST` request. For security, it sends an **OPTIONS** request
first (the "Pre-flight"). It’s basically asking: *"Are you okay with
me sending this JSON?"*

### 2. The Router Lockout (The 405 Error)

You were using the new Go 1.22 routing syntax: `mux.HandleFunc("POST
/tasks", ...)`.

* This syntax is very strict. It tells Go: **"ONLY allow POST requests
  here."**
* When the browser sent the `OPTIONS` request, the router saw that you
  didn't have a handler for `OPTIONS /tasks`.
* Instead of letting the request through to your function, the router
  automatically intercepted it and sent back a **405 Method Not
  Allowed**.

---

## How It Got Fixed

We shifted the logic from **inside** the handler to a **Middleware
Wrapper**.

### The "Shield" Approach

By wrapping your `mux` in `enableCORS(mux)`, we placed a "security
guard" in front of your router.

1. **The Interception:** When the `OPTIONS` request arrives, it hits
   your middleware **before** it reaches the router.
2. **The Approval:** The middleware sees it’s an `OPTIONS` request,
   attaches the "Allow" headers, and sends a `204 No Content` response
   immediately.
3. **The Bypass:** Because the middleware handles the `OPTIONS`
   request and returns, the strict router never sees it, so the router
   never has the chance to throw a 405 error.
4. **The Follow-through:** The browser receives the approval, then
   sends the real `POST` request. The middleware sees it's a `POST`,
   ignores the `if r.Method == "OPTIONS"` block, and passes it to the
   router, which then successfully matches `POST /tasks`.

---

## Summary of Why it Works

| Problem          | Fix                                | Why?                                          |
|------------------|------------------------------------|-----------------------------------------------|
| CORS Block       | Added `Access-Control` headers.    | Tells the browser the Frontend is trusted.    |
| Pre-flight Error | Handled `OPTIONS` method.          | Browser requires an "OK" before sending JSON. |
| 405 Error        | Used Middleware (outside the mux). | Prevents the strict Go 1.22 router            |
|                  |                                    | from rejecting unknown methods.               |

### Better Architecture

This is much better than fixing it inside the `createTask` function
because now, if you add `DELETE /tasks` or `PUT /tasks`, they are
**automatically** protected and enabled for CORS without you having to
write the headers again!

Since you have the full cycle working (Create, Read, and UI
interaction), would you like to see how to handle the **Delete**
button in your template using this same middleware?
