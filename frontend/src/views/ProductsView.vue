<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuth } from "@/stores/useAuth"
import { useProduct } from "@/stores/useProduct"

import {
	useNotification,
	NForm,
	NFormItem,
	NImage,
	NInput,
	NPopconfirm,
	NInputNumber,
	NEllipsis,
	NSelect,
	NSpace,
	NButton,
	NIcon,
	NGrid,
	NGridItem,
	NCard,
	NModal,
} from "naive-ui"

import {
	AddCircleOutline as CreateIcon,
	TrashOutline as DeleteIcon,
	SettingsOutline as UpdateIcon,
} from '@vicons/ionicons5'
const createProductFormRef = ref<Object>(null);
const updateProductFormRef = ref<Object>(null);

const productRules = {
	name: {
		message: "Please input the recipe name",
		required: true,
	},
	image_url: {
		message: "Please input a link to the image",
		required: true,
	},
	serving: {
		type: {
			message: "Please select a valid serving type",
			required: true,
		},
		size: {
			message: "Please select serving size",
			required: true,
		},
		calories: {
			message: "Please specify a single serving calories count",
			required: true,
		},
	},
	description: {
		message: "Please describe the recipe",
		required: true,
	},
}

const createdProduct = ref<Product>({});
const selectedProduct = ref<Product>({});
const createModalVisible = ref(false);
const updateModalVisible = ref(false);

const notification = useNotification();

useProduct().retrieveProducts();

const servingOptions = [
	{
		label: "Units",
		value: "units"
	},
	{
		label: "Grams",
		value: "grams"
	},
	{
		label: "Milliliters",
		value: "milliliters"
	}
]

const deleteProduct = (id: string) => {
	useProduct().deleteProduct(id).then(resp => {
		if (resp.error) {
			notification["error"]({
				content: resp.error,
				duration: 2500,
				keepAliveOnHover: true
			})
			return;
		} 

		notification["success"]({
			content: "Succesfully deleted a product",
			duration: 2500,
			keepAliveOnHover: true
		})
	});
}

const showCreateModal = () => {
	createdProduct.value = {
		serving: {}
	};
	createModalVisible.value = true;
}

const showUpdateModal = (product: Product) => {
	selectedProduct.value = JSON.parse(JSON.stringify(product));
	updateModalVisible.value = true;
}

const createProduct = () => {
	createProductFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		useProduct().createProduct(createdProduct.value).then(resp => {
			if (resp.error) {
				notification["error"]({
					content: resp.error,
					duration: 2500,
					keepAliveOnHover: true
				})
				return;
			} 

			notification["success"]({
				content: "Succesfully created a product",
				duration: 2500,
				keepAliveOnHover: true
			})

			createModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const updateProduct = () => {
	updateProductFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		useProduct().updateProduct(selectedProduct.value).then(resp => {
			if (resp.error) {
				notification["error"]({
					content: resp.error,
					duration: 2500,
					keepAliveOnHover: true
				})
				return;
			} 

			notification["success"]({
				content: "Succesfully updated a product",
				duration: 2500,
				keepAliveOnHover: true
			})

			updateModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const searchQuery = ref<string>("");

const searchedProducts = computed(() => {
	return useProduct().products.filter((product) => {
		return (
			product.name
				.toLowerCase()
				.indexOf(searchQuery.value.toLowerCase()) != -1
		);
	});
});
</script>

<template>
	<n-space justify="space-between" align="center">
		<div class="page-title">Products</div>
		<n-space justify="space-between" align="center">
			<n-input v-model:value="searchQuery" type="text" placeholder="Search" />
			<n-button v-if="useAuth().authorized && useAuth().user.admin" :focusable="false" @click="showCreateModal" text>
				<template #icon>
					<n-icon size="20px" color="#36ad6a">
						<create-icon />
					</n-icon>
				</template>
				Create
			</n-button>
		</n-space>
	</n-space>
	<n-grid cols="3 800:4 1000:5 1200:6" x-gap="20px" y-gap="20px">
		<n-grid-item v-for="(product, index) in searchedProducts" :key="index" style="display: flex">
			<n-card :title="product.name" hoverable embedded>
				<template #cover>
					<img 
						:src="product.image_url == '' ? 'https://media.istockphoto.com/vectors/thumbnail-image-vector-graphic-vector-id1147544807?k=20&m=1147544807&s=612x612&w=0&h=pBhz1dkwsCMq37Udtp9sfxbjaMl27JUapoyYpQm0anc=' : product.image_url"
					>
				</template>
				<n-ellipsis line-clamp="3">
					{{ product.description }}
				</n-ellipsis>
				<div class="nutrition-description">
					{{ product.serving.size }}
					{{ product.serving.type }} - 
					{{ product.serving.calories }}kcal
				</div>
				<template v-if="useAuth().authorized && useAuth().user.admin" #footer>
					<n-space justify="space-between" align="center">
						<n-button :focusable="false" @click="showUpdateModal(product)" text>
							<template #icon>
								<n-icon size="20px" color="#36ad6a">
									<update-icon />
								</n-icon>
							</template>
							Update
						</n-button>
						<n-popconfirm width="310px" placement="top-end" :positive-text="null" negative-text="Confirm" @negative-click="deleteProduct(product.id)">
							<template #trigger>
								<n-button type="error" :focusable="false" text>
									<template #icon>
										<n-icon size="20px" color="#d03050">
											<delete-icon />
										</n-icon>
									</template>
									Delete
								</n-button>
							</template>
							Are you sure you want to delete the product?
						</n-popconfirm>
					</n-space>
				</template>
			</n-card>
		</n-grid-item>
	</n-grid>
	<n-modal
		v-model:show="createModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			Create a new product
		</div>
		<div class="modal-content">
			<n-space justify="center">
				<div class="image-wrapper">
					<n-image
						width="250"
						object-fit="scale-down"
						:src="createdProduct.image_url"
						fallback-src="https://media.istockphoto.com/vectors/thumbnail-image-vector-graphic-vector-id1147544807?k=20&m=1147544807&s=612x612&w=0&h=pBhz1dkwsCMq37Udtp9sfxbjaMl27JUapoyYpQm0anc="
					/>
				</div>
			</n-space>
			<div class="divider"></div>
			<n-form ref="createProductFormRef" :model="createdProduct" :rules="productRules">
				<n-form-item label="Name" path="name" required>
					<n-input type="text" maxlength="255" show-count v-model:value="createdProduct.name" placeholder="Name of the product" />
				</n-form-item>
				<n-form-item label="Image URL" path="image_url" required>
					<n-input type="text" maxlength="1023" show-count v-model:value="createdProduct.image_url" placeholder="Image representation of the product" />
				</n-form-item>
				<n-form-item label="Serving Type" path="serving.type" required>
					<n-select v-model:value="createdProduct.serving.type" placeholder="How it is measured in the recipes" :options="servingOptions" />
				</n-form-item>
				<n-form-item label="Serving Size" path="serving.size" required>
					<n-input-number v-model:value="createdProduct.serving.size" placeholder="Most common size in the recipes" min="0.25" step="0.25" style="width: 100%"/>
				</n-form-item>
				<n-form-item label="Serving Calories" path="serving.calories" required>
					<n-input-number v-model:value="createdProduct.serving.calories" placeholder="Calories per serving size" :show-button="false" min="0" style="width: 100%">
						<template #suffix>
							kcal
						</template>
					</n-input-number>
				</n-form-item>
				<n-form-item label="Description" path="description">
					<n-input type="textarea" maxlength="1023" show-count v-model:value="createdProduct.description" placeholder="Description of the product" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button type="primary" @click="createProduct">
				Create
			</n-button>
		</template>
	</n-modal>
	<n-modal
		v-model:show="updateModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			Update {{ selectedProduct.name }}
		</div>
		<div class="modal-content">
			<n-space justify="center">
				<div class="image-wrapper">
					<n-image
						width="250"
						object-fit="scale-down"
						:src="selectedProduct.image_url"
						fallback-src="https://media.istockphoto.com/vectors/thumbnail-image-vector-graphic-vector-id1147544807?k=20&m=1147544807&s=612x612&w=0&h=pBhz1dkwsCMq37Udtp9sfxbjaMl27JUapoyYpQm0anc="
					/>
				</div>
			</n-space>
			<div class="divider"></div>
			<n-form ref="updateProductFormRef" :model="selectedProduct" :rules="productRules">
				<n-form-item label="Image URL" path="image_url" required>
					<n-input type="text" maxlength="1023" show-count v-model:value="selectedProduct.image_url" placeholder="Image representation of the product" />
				</n-form-item>
				<n-form-item label="Serving Type" path="serving.type" required>
					<n-select v-model:value="selectedProduct.serving.type" placeholder="How it is measured in the recipes" :options="servingOptions" />
				</n-form-item>
				<n-form-item label="Serving Size" path="serving.size" required>
					<n-input-number v-model:value="selectedProduct.serving.size" placeholder="Most common size in the recipes" min="0.25" step="0.25" style="width: 100%"/>
				</n-form-item>
				<n-form-item label="Serving Calories" path="serving.calories" required>
					<n-input-number v-model:value="selectedProduct.serving.calories" placeholder="Calories per serving size" :show-button="false" min="0" style="width: 100%">
						<template #suffix>
							kcal
						</template>
					</n-input-number>
				</n-form-item>
				<n-form-item label="Description" path="description">
					<n-input type="textarea" maxlength="1023" show-count v-model:value="selectedProduct.description" placeholder="Description of the product" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button type="primary" @click="updateProduct">
				Update
			</n-button>
		</template>
	</n-modal>
</template>

<style lang="scss" scoped>
@use "../assets/main.scss";

.modal-title {
	color: main.$green;
	font-weight: 600;
	font-size: 28px;
}

.modal-content {
	.image-wrapper {
		margin: 20px 0;
		border: 1px solid main.$grey;
	}
}

.nutrition-description {
	color: main.$green;
	margin-bottom: 5px;
}

.divider {
	margin-bottom: 20px;
}
</style>
