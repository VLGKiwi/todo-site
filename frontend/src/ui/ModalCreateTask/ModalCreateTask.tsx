'use client'

import React from 'react'
import { FC } from 'react'
import { ModalCreateTaskProps } from './ModalCreateTask.types'
import { useCreateTodo } from "@/hooks/useCreateTodo";

export const ModalCreateTask: FC<ModalCreateTaskProps> = ({isOpen}) => {

	const { mutate, isPending } = useCreateTodo()

	console.log(isOpen)

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		const form = e.currentTarget;
		const fd = new FormData(form);
		const title = String(fd.get("title") ?? "Нет названия");
		const description = String(fd.get("description") ?? "Нет описания");

		mutate({title, description})
	}

	return (
		<div>
			<form onSubmit={handleSubmit}>
				<label htmlFor="">
					Заголовок
					<input type="text" name='title' />
				</label>
				<label htmlFor="">
					Описание
					<input type="text" name='description' />
				</label>
				{
					isPending ?
						<p>Отправка</p> :
						<button type='submit'>Сохранить</button>
				}
			</form>
		</div>
	)
}
