export type CreateTodoResponse = {
	id: number
}

export type GetAllTodosResponse = Todo[]

export type Todo = {
	id: number
	title: string
	description: string
	completed: boolean
}

export type UpdateTodoRequest = {
	message: string
}
