import { ref } from "vue"
import { defineStore } from "pinia"
import { useAuth } from "@/stores/useAuth"
import axios from 'axios'

export const useProduct = defineStore("product", () => {
	interface Product {
		id: string
		name: string
		description: string
		image_url: string
		serving: {
			type: string
			size: number
			calories: number
		}
	}

	const products = ref<Array<Product>>([]);

	function retrieveProduct(id: string): Product {
		return axios
			.get("/api/products/" + id)
			.then((response) => {
				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function retrieveProducts() {
		return axios
			.get("/api/products")
			.then((response) => {
				if (response.status == 200) {
					products.value = response.data;
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function createProduct(product: Product) {
		return axios
			.post("/api/products", product, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					retrieveProducts();
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function updateProduct(product: Product) {
		return axios
			.patch("/api/products/" + product.id, product, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					retrieveProducts();
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function deleteProduct(id: string) {
		return axios
			.delete("/api/products/" + id, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 204) {
					retrieveProducts();
				}

				return {}
			})
			.catch((error) => {
				if (error.response.status == 409) {
					return {
						error: "Product is in use.",
					}
				}

				return {
					error: error.message,
				}
			})
	}

	return { products, retrieveProduct, retrieveProducts, createProduct, updateProduct, deleteProduct }
})
