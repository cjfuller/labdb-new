package routes

import "github.com/gin-gonic/gin"

func InstallAll(r *gin.Engine) {
	installplasmid(r)
	installplasmids(r)
	installoligo(r)
	installoligos(r)
	installline(r)
	installlines(r)
	installsample(r)
	installsamples(r)
	installbacterium(r)
	installbacteria(r)
	installyeaststrain(r)
	installyeaststrains(r)
	installuser(r)
	installusers(r)
	installantibody(r)
	installantibodies(r)
	installrnai_clone(r)
	installrnai_clones(r)
	installseq_lib(r)
	installseq_libs(r)

}
func installplasmid(r *gin.Engine) {
	r.GET("/plasmid/:id/next", nextRoute("plasmid"))
	r.GET("/plasmid/:id/previous", previousRoute("plasmid"))
}
func installplasmids(r *gin.Engine) {
	r.GET("/plasmids/:id/next", nextRoute("plasmids"))
	r.GET("/plasmids/:id/previous", previousRoute("plasmids"))
}
func installoligo(r *gin.Engine) {
	r.GET("/oligo/:id/next", nextRoute("oligo"))
	r.GET("/oligo/:id/previous", previousRoute("oligo"))
}
func installoligos(r *gin.Engine) {
	r.GET("/oligos/:id/next", nextRoute("oligos"))
	r.GET("/oligos/:id/previous", previousRoute("oligos"))
}
func installline(r *gin.Engine) {
	r.GET("/line/:id/next", nextRoute("line"))
	r.GET("/line/:id/previous", previousRoute("line"))
}
func installlines(r *gin.Engine) {
	r.GET("/lines/:id/next", nextRoute("lines"))
	r.GET("/lines/:id/previous", previousRoute("lines"))
}
func installsample(r *gin.Engine) {
	r.GET("/sample/:id/next", nextRoute("sample"))
	r.GET("/sample/:id/previous", previousRoute("sample"))
}
func installsamples(r *gin.Engine) {
	r.GET("/samples/:id/next", nextRoute("samples"))
	r.GET("/samples/:id/previous", previousRoute("samples"))
}
func installbacterium(r *gin.Engine) {
	r.GET("/bacterium/:id/next", nextRoute("bacterium"))
	r.GET("/bacterium/:id/previous", previousRoute("bacterium"))
}
func installbacteria(r *gin.Engine) {
	r.GET("/bacteria/:id/next", nextRoute("bacteria"))
	r.GET("/bacteria/:id/previous", previousRoute("bacteria"))
}
func installyeaststrain(r *gin.Engine) {
	r.GET("/yeaststrain/:id/next", nextRoute("yeaststrain"))
	r.GET("/yeaststrain/:id/previous", previousRoute("yeaststrain"))
}
func installyeaststrains(r *gin.Engine) {
	r.GET("/yeaststrains/:id/next", nextRoute("yeaststrains"))
	r.GET("/yeaststrains/:id/previous", previousRoute("yeaststrains"))
}
func installuser(r *gin.Engine) {
	r.GET("/user/:id/next", nextRoute("user"))
	r.GET("/user/:id/previous", previousRoute("user"))
}
func installusers(r *gin.Engine) {
	r.GET("/users/:id/next", nextRoute("users"))
	r.GET("/users/:id/previous", previousRoute("users"))
}
func installantibody(r *gin.Engine) {
	r.GET("/antibody/:id/next", nextRoute("antibody"))
	r.GET("/antibody/:id/previous", previousRoute("antibody"))
}
func installantibodies(r *gin.Engine) {
	r.GET("/antibodies/:id/next", nextRoute("antibodies"))
	r.GET("/antibodies/:id/previous", previousRoute("antibodies"))
}
func installrnai_clone(r *gin.Engine) {
	r.GET("/rnai_clone/:id/next", nextRoute("rnai_clone"))
	r.GET("/rnai_clone/:id/previous", previousRoute("rnai_clone"))
}
func installrnai_clones(r *gin.Engine) {
	r.GET("/rnai_clones/:id/next", nextRoute("rnai_clones"))
	r.GET("/rnai_clones/:id/previous", previousRoute("rnai_clones"))
}
func installseq_lib(r *gin.Engine) {
	r.GET("/seq_lib/:id/next", nextRoute("seq_lib"))
	r.GET("/seq_lib/:id/previous", previousRoute("seq_lib"))
}
func installseq_libs(r *gin.Engine) {
	r.GET("/seq_libs/:id/next", nextRoute("seq_libs"))
	r.GET("/seq_libs/:id/previous", previousRoute("seq_libs"))
}
