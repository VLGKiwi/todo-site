import { api } from "@/lib/api";
import { GetAllTodos } from "@/types/todo";
import { useQuery } from "@tanstack/react-query";

export function useTodos() {
	return useQuery({
		queryKey: ['todos'],
		queryFn: async () => {
			const { data } = await api.get<GetAllTodos>('/todos');
			return data;
		}
	})
}
