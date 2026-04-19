'use client'

import React from 'react'
import { FC } from 'react'
import { ModalCreateTaskProps } from './ModalCreateTask.types'
import { useCreateTodo } from "@/hooks/useCreateTodo";
import styles from './ModalCreateTask.module.scss'
import { Button } from '../Button/Button';
import Image from 'next/image';
import { useCreateTaskModalStore } from '@/stores/createTaskModalStore';

export const ModalCreateTask: FC<ModalCreateTaskProps> = ({isOpen}) => {

	const { mutate, isPending } = useCreateTodo()
	const close = useCreateTaskModalStore((s) => s.close)

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		const form = e.currentTarget;
		const fd = new FormData(form);
		const title = String(fd.get("title") ?? "Нет названия");
		const description = String(fd.get("description") ?? "Нет описания");

		mutate({ title, description },
			{
				onSuccess: () => {
					close();
					form.reset()
				},
				onError: (error) => {
					console.error(error) // Заглушка на обработку ошибки
				}
			}
		)
	}

	const visible = !isOpen ? styles.visible : ''

	return (
		<div
			onClick={close}
			className={`${styles.container} ${visible}`}
		>
			<Image
				src={'/icons/cross.svg'}
				width={32}
				height={32}
				alt='Закрытие модалки'
				className={styles.close}
				onClick={close}
			/>
			<form
				className={styles.modal}
				onSubmit={handleSubmit}

				onClick={(e) => e.stopPropagation()}
			>
				<h2 className={styles.title}>
					Создание задачи
				</h2>
				<label className={styles.line} htmlFor="">
					Заголовок
					<input type="text" name='title' />
				</label>
				<label className={styles.line} htmlFor="">
					Описание
					<input type="text" name='description' />
				</label>
				{
					isPending ?
						<p>Сохранение</p> :
						<Button type='submit'>Сохранить</Button>
				}
			</form>
		</div>
	)
}
