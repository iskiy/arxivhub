package api

//func (rest *Rest) auth(next http.HandlerFunc) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token := r.Header.Get("Authorization")
//
//		fields := strings.Fields(token)
//		if len(fields) < 2 {
//			err := errors.New("invalid authorization header format")
//			rest.sendError(w, http.StatusUnauthorized, err)
//			return
//		}
//
//		authorizationType := strings.ToLower(fields[0])
//		if authorizationType != "bearer" {
//			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
//			rest.sendError(w, http.StatusUnauthorized, err)
//			return
//		}
//
//		accessToken := fields[1]
//
//		payload, err := rest.tokenMaker.VerifyToken(accessToken)
//		if err != nil {
//			rest.sendError(w, http.StatusUnauthorized, err)
//			return
//		}
//		fmt.Println(payload)
//		payloadContext := context.WithValue(r.Context(), "payload", payload)
//		fmt.Println(payloadContext.Value("payload"))
//
//		next.ServeHTTP(w, r.WithContext(payloadContext))
//	})
//}
