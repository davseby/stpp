import { createRouter, createWebHistory } from "vue-router"
import { useAuth } from "@/stores/useAuth"
import ProductsView from "@/views/ProductsView.vue"
import RecipesView from "@/views/RecipesView.vue"
import PlansView from "@/views/PlansView.vue"
import UsersView from "@/views/UsersView.vue"

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/",
			redirect: "/products",
		},
		{
			path: "/products",
			name: "products",
			component: ProductsView,
		},
		{
			path: "/recipes",
			name: "recipes",
			component: RecipesView,
		},
		{
			path: "/plans",
			name: "plans",
			component: PlansView,
		},
		{
			path: "/users",
			name: "users",
			component: UsersView,
			beforeEnter: (to, from) => {
				return useAuth().authorized && useAuth().user.admin
			},
		},
	],
})

export default router
