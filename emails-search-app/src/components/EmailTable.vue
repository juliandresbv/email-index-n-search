<template>
  <div class="mx-7">
    <div class="mb-3">
      <h2 class="text-xl font-semibold leading-tight">Emails</h2>
    </div>
    <div class="border inline-block shadow rounded-lg max-h-[80vh]">
      <div class="overflow-y-auto max-h-[70vh]">
        <table class="w-full table-fixed">
          <tbody v-if="emails?.length > 0">
            <tr
              class="cursor-pointer bg-white hover:bg-gray-200"
              v-for="email in emails"
              :key="email.id"
              @click="selectEmail(email)"
            >
              <td class="px-5 py-5 border-b border-gray-200 text-sm">
                <div class="flex-col">
                  <div class="flex justify-between">
                    <p class="w-2/3 text-gray-900 whitespace-no-wrap truncate text-left font-bold">
                      {{ email.subject ? email.subject : '(No subject)' }}
                    </p>
                    <p class="w-[20vh] text-gray-500 whitespace-no-wrap truncate text-right">
                      {{ email.date }}
                    </p>
                  </div>
                  <div class="mt-4">
                    <p class="text-gray-900 whitespace-no-wrap truncate text-left">
                      {{ email.body }}
                    </p>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr class="bg-white" v-for="n in 10" :key="'empty-row' + n">
              <td v-if="n == 5" class="px-5 py-5 border-b border-gray-200 text-sm">
                <p class="text-gray-500 text-center">No data</p>
              </td>
              <td v-else class="px-5 py-5 border-b border-gray-200 text-sm">
                <p class="text-gray-500 text-center">&nbsp;</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex items-center justify-between bg-white">
        <div class="w-1/2 px-5 py-5 flex flex-col xs:flex-row items-center xs:justify-between">
          <span class="text-xs xs:text-sm text-gray-900 font-semibold">
            Current limit: {{ limit }}
          </span>
          <div class="inline-flex mt-2 xs:mt-0">
            <select
              class="text-sm text-gray-800 py-2 px-4 mx-2 rounded border"
              name="email-search-limit"
              id="email-search-limit"
              @change="changeLimit"
            >
              <option value="10">10</option>
              <option value="20">20</option>
              <option value="50">50</option>
            </select>
          </div>
        </div>
        <div class="w-1/2 px-5 py-5 flex flex-col xs:flex-row items-center xs:justify-between">
          <span class="text-xs xs:text-sm text-gray-900 font-semibold">
            Showing {{ page }} to {{ maxPages }} of {{ hits }} entries
          </span>
          <div class="inline-flex mt-2 xs:mt-0">
            <button
              class="w-1/4 text-sm bg-gray-300 text-gray-800 py-2 px-4 mx-2 rounded font-semibold"
              @click="prevPage"
              :disabled="page <= 1"
              :class="{ 'opacity-50': page <= 1 }"
            >
              &lt;
            </button>
            <input
              class="w-1/2 text-sm text-gray-800 py-2 px-4 mx-2 rounded border"
              id="page-input"
              name="page-input"
              type="number"
              v-model="page"
              :disabled="hits <= 0"
              :min="1"
              :max="hits > 0 ? maxPages : 0"
              @change="changePage"
            />
            <button
              class="w-1/4 text-sm bg-gray-300 text-gray-800 py-2 px-4 rounded mx-2 font-semibold"
              @click="nextPage"
              :disabled="page >= maxPages"
              :class="{ 'opacity-50': page >= maxPages }"
            >
              &gt;
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'

import { makeRequest } from '../services'
import type { Email, EmailsSearchResponse } from '../types'
import { useEmailsStore, useComponentsRefsStore } from '../stores'

const emailsStore = useEmailsStore()
const componentsRefsStore = useComponentsRefsStore()

const { hits } = storeToRefs(emailsStore)
const { emails } = storeToRefs(emailsStore)

const { page } = storeToRefs(componentsRefsStore)
const { limit } = storeToRefs(componentsRefsStore)
const { searchTerm } = storeToRefs(componentsRefsStore)

const maxPages = ref(0)

const debounceTimeout = 500

let timeoutId: NodeJS.Timeout | null = null

watch(hits, (newHits) => {
  if (newHits <= 0) {
    componentsRefsStore.setPage(1)
    maxPages.value = 0
  } else {
    componentsRefsStore.setPage(1)
    maxPages.value = Math.ceil(newHits / limit.value)
  }
})

const selectEmail = (email: Email) => {
  emailsStore.setSelectedEmail(email)
}

const prevPage = async () => {
  if (page.value - 1 <= 0) {
    return
  }

  if (searchTerm.value.length <= 0) {
    return
  }

  componentsRefsStore.setPage(page.value - 1)

  if (timeoutId) {
    clearTimeout(timeoutId)
  }

  timeoutId = setTimeout(async () => {
    const searchEmailsReq = await makeRequest<EmailsSearchResponse>('POST', '/emails/search', {
      term: searchTerm.value,
      page: page.value,
      limit: limit.value
    })

    if (!searchEmailsReq) {
      emailsStore.$reset()

      return
    }

    const { emails, hits } = searchEmailsReq.data

    if (!emails || emails?.length <= 0 || !hits || hits <= 0) {
      emailsStore.$reset()

      return
    }

    emailsStore.setEmails(emails)
    emailsStore.setHits(hits)
  }, debounceTimeout)
}

const nextPage = async () => {
  if (page.value + 1 > maxPages.value) {
    return
  }

  if (searchTerm.value.length <= 0) {
    return
  }

  componentsRefsStore.setPage(page.value + 1)

  if (timeoutId) {
    clearTimeout(timeoutId)
  }

  timeoutId = setTimeout(async () => {
    const searchEmailsReq = await makeRequest<EmailsSearchResponse>('POST', '/emails/search', {
      term: searchTerm.value,
      page: page.value,
      limit: limit.value
    })

    if (!searchEmailsReq) {
      emailsStore.$reset()

      return
    }

    const { hits, emails } = searchEmailsReq.data

    if (!emails || emails.length <= 0 || !hits || hits <= 0) {
      emailsStore.$reset()

      return
    }

    emailsStore.setEmails(emails)
    emailsStore.setHits(hits)
  }, debounceTimeout)
}

const changeLimit = async (event: any) => {
  const newLimit = Number(event.target.value)

  componentsRefsStore.setLimit(newLimit)
  componentsRefsStore.setPage(1)
  maxPages.value = Math.ceil(hits.value / limit.value)

  if (searchTerm.value.length <= 0) {
    return
  }

  if (timeoutId) {
    clearTimeout(timeoutId)
  }

  timeoutId = setTimeout(async () => {
    const searchEmailsReq = await makeRequest<EmailsSearchResponse>('POST', '/emails/search', {
      term: searchTerm.value,
      page: page.value,
      limit: limit.value
    })

    if (!searchEmailsReq) {
      emailsStore.$reset()

      return
    }

    const { emails, hits } = searchEmailsReq.data

    if (!emails || emails?.length <= 0 || !hits || hits <= 0) {
      emailsStore.$reset()

      return
    }

    emailsStore.setEmails(emails)
    emailsStore.setHits(hits)
  }, debounceTimeout)
}

const changePage = (event: any) => {
  let newPage = Number(event.target.value)

  if (newPage <= 0) {
    newPage = 1
  }

  if (newPage > maxPages.value) {
    newPage = maxPages.value
  }

  if (searchTerm.value.length <= 0) {
    return
  }

  componentsRefsStore.setPage(newPage)

  if (timeoutId) {
    clearTimeout(timeoutId)
  }

  timeoutId = setTimeout(async () => {
    const searchEmailsReq = await makeRequest<EmailsSearchResponse>('POST', '/emails/search', {
      term: searchTerm.value,
      page: page.value,
      limit: limit.value
    })

    if (!searchEmailsReq) {
      emailsStore.$reset()

      return
    }

    const { emails, hits } = searchEmailsReq.data

    if (!emails || emails?.length <= 0 || !hits || hits <= 0) {
      emailsStore.$reset()

      return
    }

    emailsStore.setEmails(emails)
    emailsStore.setHits(hits)
  }, debounceTimeout)
}
</script>
