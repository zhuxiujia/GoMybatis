#ifdef GOARCH_arm
#define LR R14
#endif

#ifdef GOARCH_amd64
#define	get_tls(r)	MOVQ TLS, r
#define	g(r)	0(r)(TLS*1)
#endif

#ifdef GOARCH_amd64p32
#define	get_tls(r)	MOVL TLS, r
#define	g(r)	0(r)(TLS*1)
#endif

#ifdef GOARCH_386
#define	get_tls(r)	MOVL TLS, r
#define	g(r)	0(r)(TLS*1)
#endif