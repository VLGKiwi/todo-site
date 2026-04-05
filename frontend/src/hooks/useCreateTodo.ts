import { api } from "@/lib/api";
import { CreateTodoResponse, Todo } from "@/types/todo";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useCreateTodo() {
	const queryClient = useQueryClient()

	return useMutation({
		mutationFn: async (body: Omit<Todo, 'id' | 'completed'>) => {
			const { data } = await api.post<CreateTodoResponse>('/todos', body);

			return data
		},

		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['todos'] })
		}
	})
}
