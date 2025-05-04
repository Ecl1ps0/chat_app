<script setup lang="ts">
import { ref, onMounted, watch, nextTick, computed } from 'vue';
import { useAuth } from '@/shared/stores/auth.store';
import { getAvailableUsers, updateUser } from '../api/users.api';
import { MessageSquareIcon, SendIcon, ImageIcon, XIcon, MenuIcon, LogOutIcon, SearchIcon } from 'lucide-vue-next';
import { router } from '@/pages/router/Router';
import { useChatSocket } from '../api/chat.api';
import { toast } from '@/components/ui/toast';
import UserProfile from '@/pages/home/components/Profile.vue';
import { Button } from '@/components/ui/button';
import { Sheet, SheetTrigger, SheetContent, SheetHeader, SheetTitle, SheetDescription, SheetFooter } from '@/components/ui/sheet';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { IUser } from '@/entities/user.entity';
import { jwtDecode } from 'jwt-decode';
import Input from '@/components/ui/input/Input.vue';
import { IMessage } from '@/entities/message.entity';
import ContextMenu from '../components/ContextMenu.vue';
import { uploadImages } from '../api/image.api';
import AudioRecorder from '../components/AudioRecorder.vue';
import { uploadAudio } from '../api/audio.api';

const https_domain = import.meta.env.VITE_DOMAIN_HTTPS;

const users = ref<IUser[]>([]);
const newMessage = ref('');
const selectedFiles = ref<File[]>([]);
const isUploading = ref(false);
const isProfileOpen = ref(false);
const selectedUserProfile = ref<IUser | null>(null);
const contextMenu = ref({ show: false, x: 0, y: 0, isCurrentUser: false, isAudioMessage: false });
const selectedMessage = ref<IMessage | null >(null);
const isUpdateModalOpen = ref(false);
const updatedMessageContent = ref('');
const isDeleteModalOpen = ref(false);
const deleteForBoth = ref(false);

const { token, isAuthorized, isExpired, setToken, logout } = useAuth();
const currentUser: IUser = jwtDecode<{ exp: number; user: IUser }>(token!).user;
const { messages, selectedUser, selectUser, sendMessage, updateMessage, deleteMessage, closeConnection } = useChatSocket(currentUser.id);

const handleSelectUser = (user: IUser) => {
  if (selectedUser.value) {
    closeConnection();
  }
  selectUser(token!, user);
};

const handleFileSelect = (event: Event) => {
  const fileList = (event.target as HTMLInputElement).files;
  if (fileList) {
    selectedFiles.value = Array.from(fileList);
  }
};

const removeFile = (index: number) => {
  selectedFiles.value.splice(index, 1);
};

const handleLogout = async () => {
  try {
    logout();
    router.push('/auth');
  } catch (e) {
    toast({
      title: 'Error',
      description: `Failed to logout: ${(e as Error).message}`
    });
  }
};

const handleSubmit = async () => {
  if (newMessage.value.trim() === '' && selectedFiles.value.length === 0) return;

  isUploading.value = true;
  try {
    let imageIds: string[] = [];
    if (selectedFiles.value.length > 0) {
      const formData = new FormData();
      selectedFiles.value.forEach((file: File) => {
        formData.append("images", file); 
      });

      imageIds = await uploadImages(token!, formData);
    }

    await sendMessage(newMessage.value, imageIds, undefined);
    newMessage.value = '';
    selectedFiles.value = [];
  } catch (e) {
    toast({
      title: (e as Error).name,
      description: `Error uploading files: ${(e as Error).message}`
    });
  } finally {
    isUploading.value = false;
  }
};

const handleUpdateMessage = async () => {
  if (!selectedMessage.value || !updatedMessageContent.value.trim()) return;

  try {
    await updateMessage(selectedMessage.value.id, updatedMessageContent.value);
    
    // Update the message in the local state
    const index = messages.value.findIndex(m => m.id === selectedMessage.value!.id);
    if (index !== -1) {
      messages.value[index].text_content = updatedMessageContent.value;
    }

    toast({
      title: 'Success',
      description: 'Message updated successfully'
    });
    closeUpdateModal();
  } catch (error) {
    toast({
      title: 'Error',
      description: `Failed to update message: ${(error as Error).message}`
    });
  }
};

const handleDeleteMessage = async () => {
  if (!selectedMessage.value) return;

  console.log(selectedMessage.value)

  try {
    await deleteMessage(selectedMessage.value.id, deleteForBoth.value, selectedUser.value?.id);
    
    const index = messages.value.findIndex(m => m.id === selectedMessage.value!.id);
    if (index !== -1) {
      messages.value.splice(index, 1);
    }

    toast({
      title: 'Success',
      description: 'Message deleted successfully'
    });
    closeDeleteModal();
  } catch (error) {
    toast({
      title: 'Error',
      description: `Failed to delete message: ${(error as Error).message}`
    });
  }
};

const handleAudioRecording = async (audioBlob: Blob) => {
  isUploading.value = true;
  try {
    const formData = new FormData();
    formData.append("audio", audioBlob, "audio_message.wav");

    // Assuming you have an API endpoint for uploading audio
    const audioId = await uploadAudio(token!, formData);

    // Send the message with the audio ID
    await sendMessage(undefined, undefined, audioId.audio_id);
  } catch (e) {
    toast({
      title: (e as Error).name,
      description: `Error uploading audio: ${(e as Error).message}`
    });
  } finally {
    isUploading.value = false;
  }
};

const handleUpdateUser = async (formData: FormData) => {
  try {
    const newToken = await updateUser(token!, formData);
    setToken(newToken.access_token);
    toast({
      title: 'Success',
      description: 'Profile updated successfully'
    });
  } catch (e) {
    console.log(e)
    toast({
      title: 'Error',
      description: `Failed to update profile: ${(e as Error).message}`
    });
  }
};

const showUserProfile = (user: IUser) => {
  selectedUserProfile.value = user;
  isProfileOpen.value = true;
};

onMounted(async () => {
  if (!isAuthorized() || isExpired()) {
    router.push("/auth");
  } else {
    users.value = await getAvailableUsers(token!);
  }
});

// Code for context menu
const handleContextMenu = (event: MouseEvent, message: IMessage) => {
  event.preventDefault();
  contextMenu.value = { 
    show: true, 
    x: event.clientX, 
    y: event.clientY,
    isCurrentUser: message.sender === currentUser.id,
    isAudioMessage: !!message.audio_content
  };
  selectedMessage.value = message;
};

const closeContextMenu = () => {
  contextMenu.value.show = false;
};

const openUpdateModal = () => {
  isUpdateModalOpen.value = true;
  updatedMessageContent.value = selectedMessage.value!.text_content;
  closeContextMenu();
};

const closeUpdateModal = () => {
  isUpdateModalOpen.value = false;
  selectedMessage.value = null;
  updatedMessageContent.value = '';
};

const openDeleteModal = () => {
  isDeleteModalOpen.value = true;
  deleteForBoth.value = false;
  closeContextMenu();
};

const closeDeleteModal = () => {
  isDeleteModalOpen.value = false;
  selectedMessage.value = null;
};

// Code for searching
const searchQuery = ref('');
const isSearching = ref(false);

const filteredMessages = computed(() => {
  if (!searchQuery.value) return messages.value;
  const query = searchQuery.value.toLowerCase();
  return messages.value.filter(message => 
    message.text_content.toLowerCase().includes(query)
  );
});

const toggleSearch = () => {
  isSearching.value = !isSearching.value;
  if (!isSearching.value) {
    searchQuery.value = '';
  } else {
    nextTick(() => {
      const searchInput = document.getElementById('message-search');
      if (searchInput) {
        searchInput.focus();
      }
    });
  }
};

// Code for auto-scrolling
const messageContainer = ref<HTMLElement | null>(null);

const scrollToBottom = () => {
  if (messageContainer.value && !searchQuery.value) {
    messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
  }
};

watch(messages, () => {
  nextTick(() => {
    scrollToBottom();
  });
}, { deep: true });

const getDelay = (index: number) => `${index * 50}ms`;

const beforeEnter = (el: Element) => {
  const htmlEl = el as HTMLElement;
  htmlEl.style.opacity = '0';
  htmlEl.style.transform = 'translateY(100px)';
};

const enter = (el: Element, done: () => void) => {
  const htmlEl = el as HTMLElement;
  const delay = getDelay(Number(htmlEl.dataset.index));
  htmlEl.style.transition = `all 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94) ${delay}`;
  htmlEl.style.opacity = '1';
  htmlEl.style.transform = 'translateY(0)';
  htmlEl.addEventListener('transitionend', done);
};

const leave = (el: Element, done: () => void) => {
  const htmlEl = el as HTMLElement;
  htmlEl.style.transition = 'all 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94)';
  htmlEl.style.opacity = '0';
  htmlEl.style.transform = 'translateY(100px)';
  htmlEl.addEventListener('transitionend', done);
};
</script>

<style scoped>
.message-enter-active,
.message-leave-active {
  transition: all 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
}
.message-enter-from,
.message-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
.message-move {
  transition: all 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
}
</style>

<template>
  <div class="flex h-screen bg-background text-foreground">
    <!-- User list sidebar -->
    <aside class="w-64 border-r bg-muted overflow-y-auto">
      <div class="p-4 border-b flex justify-between items-center">
        <h2 class="text-lg font-semibold">Contacts</h2>
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon">
              <MenuIcon class="h-5 w-5" />
              <span class="sr-only">Open menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>Your Profile</SheetTitle>
              <SheetDescription>View and edit your profile information</SheetDescription>
            </SheetHeader>
            <UserProfile 
              :user="currentUser!" 
              :is-current-user="true" 
              @update-user="handleUpdateUser" 
            />
            <SheetFooter>
              <Button variant="destructive" @click="handleLogout" class="w-full mt-4">
                <LogOutIcon class="mr-2 h-4 w-4" />
                Logout
              </Button>
            </SheetFooter>
          </SheetContent>
        </Sheet>
      </div>
      <ul>
        <li 
          v-for="user in users" 
          :key="user.id" 
          class="flex items-center p-4 hover:bg-muted-hover cursor-pointer transition-colors"
          :class="{ 'bg-primary/10': selectedUser?.id === user.id }"
        >
          <div 
            @click="handleSelectUser(user)"
            class="flex items-center flex-1"
          >
            <div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-semibold mr-3">
              <img v-if="user.profile_picture" :src="`${https_domain}/api/image?id=${user.profile_picture}`" :alt="user.username" class="w-full h-full object-center rounded-full">
              <p v-else>{{ user.username.charAt(0) }}</p>
            </div>
            <div>
              <p class="font-medium">{{ user.username }}</p>
            </div>
          </div>
          <Button 
            variant="ghost" 
            size="icon" 
            @click="showUserProfile(user)"
          >
            <MenuIcon class="h-4 w-4" />
            <span class="sr-only">View {{ user.username }}'s profile</span>
          </Button>
        </li>
      </ul>
    </aside>

    <!-- Chat area -->
    <main class="flex-1 flex flex-col">
      <header v-if="selectedUser" class="border-b p-4 flex items-center justify-between">
        <div class="flex items-center">
          <div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-semibold mr-3">
            <img v-if="selectedUser.profile_picture" :src="`${https_domain}/api/image?id=${selectedUser.profile_picture}`" :alt="selectedUser.username" class="w-full h-full object-center rounded-full">
            <p v-else>{{ selectedUser.username.charAt(0) }}</p>
          </div>
          <div>
            <h2 class="font-semibold">{{ selectedUser.username }}</h2>
          </div>
        </div>
        <Button variant="ghost" size="icon" @click="toggleSearch">
          <SearchIcon class="h-5 w-5" />
          <span class="sr-only">Search messages</span>
        </Button>
      </header>

      <div v-if="isSearching" class="p-2 border-b">
        <Input
          id="message-search"
          v-model="searchQuery"
          type="search"
          placeholder="Search messages..."
          class="w-full"
        />
      </div>

      <div v-if="selectedUser" ref="messageContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
        <TransitionGroup 
          name="message" 
          tag="ul" 
          class="space-y-4"
          :css="false"
          @before-enter="beforeEnter"
          @enter="enter"
          @leave="leave"
        >
          <li v-for="(message, index) in filteredMessages" :key="message.id" 
              :data-index="index"
              :class="[
                'flex',
                'items-start',
                'space-x-2',
                'p-4',
                'rounded-lg',
                'max-w-[80%]',
                message.sender === currentUser!.id ? 'bg-primary text-primary-foreground ml-auto' : 'bg-muted',
                message.sender === currentUser!.id ? 'flex-row-reverse' : 'flex-row'
              ]"
              @contextmenu="handleContextMenu($event, message)"
          >
            <div :class="[
              'w-8',
              'h-8',
              'rounded-full',
              'flex',
              'items-center',
              'justify-center',
              'text-xs',
              'font-semibold',
              message.sender === currentUser!.id ? 'bg-primary-foreground text-primary' : 'bg-primary text-primary-foreground'
            ]">
              <p v-if="message.sender === currentUser!.id">You</p>
              <p v-else-if="!selectedUser.profile_picture">{{selectedUser.username.charAt(0)}}</p>
              <img v-else :src="`${https_domain}/api/image?id=${selectedUser.profile_picture}`" :alt="selectedUser.username" class="w-full h-full object-center rounded-full">
            </div>
            <div class="flex-1">
              <p>{{ message.text_content }}</p>
              <div v-if="message.image_content?.length > 0" class="mt-2 flex flex-wrap gap-4">
                <div v-for="(image, imgIndex) in message.image_content" :key="imgIndex" class="w-1/4">
                  <img :src="image" alt="Uploaded image" loading="lazy" class="rounded-md w-full h-auto"/>
                </div>
              </div>
              <div v-if="message.audio_content?.length > 0" class="mt-2">
                <audio :src="`${https_domain}/api/audio?audio_id=${message.audio_content}`" controls class="w-full"></audio>
              </div>
              <span class="text-xs text-muted-foreground mt-1 block">
                {{ new Date(message.created_at * 1000).toLocaleString("en-GB", {
                  day: '2-digit', 
                  month: '2-digit', 
                  year: 'numeric', 
                  hour: '2-digit', 
                  minute: '2-digit', 
                  second: '2-digit', 
                  hour12: false
                }) }}
              </span>
            </div>
          </li>
        </TransitionGroup>
      </div>

      <div v-if="filteredMessages.length === 0 && searchQuery" class="flex-1 flex items-center justify-center bg-muted/30">
        <div class="text-center">
          <SearchIcon class="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h2 class="text-xl font-semibold mb-2">No messages found</h2>
          <p class="text-muted-foreground">Try a different search term</p>
        </div>
      </div>

      <div v-else-if="!selectedUser" class="flex-1 flex items-center justify-center bg-muted/30">
        <div class="text-center">
          <MessageSquareIcon class="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h2 class="text-xl font-semibold mb-2">Start a conversation</h2>
          <p class="text-muted-foreground">Select a contact to begin chatting</p>
        </div>
      </div>

      <footer v-if="selectedUser" class="border-t p-4">
        <form @submit.prevent="handleSubmit" class="space-y-2">
          <div class="flex space-x-2">
            <div class="relative flex-1">
              <Input
                v-model="newMessage"
                type="text"
                placeholder="Type a message..."
                class="w-full px-3 py-2 bg-background border rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <label class="px-4 py-2 bg-muted text-muted-foreground rounded-md hover:bg-muted-hover cursor-pointer">
              <ImageIcon class="h-5 w-5" />
              <Input type="file" class="hidden" multiple @change="handleFileSelect" accept="image/*" />
            </label>
            <AudioRecorder @sendAudio="handleAudioRecording" />
            <Button
              type="submit"
              class="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary"
              :disabled="isUploading"
            >
              <SendIcon v-if="!isUploading" class="h-5 w-5" />
              <span v-else class="loading loading-spinner loading-sm"></span>
            </Button>
          </div>
          <div v-if="selectedFiles.length > 0" class="flex flex-wrap gap-2 mt-2">
            <div v-for="(file, index) in selectedFiles" :key="index" class="bg-muted rounded-md p-2 flex items-center">
              <span class="text-sm truncate max-w-[100px]">{{ file.name }}</span>
              <button @click="removeFile(index)" class="ml-2 text-muted-foreground hover:text-foreground">
                <XIcon class="h-4 w-4" />
              </button>
            </div>
          </div>
        </form>
      </footer>
    </main>

    <!-- Context Menu -->
    <ContextMenu
      :show="contextMenu.show"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :isCurrentUser="contextMenu.isCurrentUser"
      :isAudioMessage="contextMenu.isAudioMessage"
      @update="openUpdateModal"
      @delete="openDeleteModal"
      @close="closeContextMenu"
    />

    <!-- Update Message Modal -->
    <Dialog v-model:open="isUpdateModalOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Update Message</DialogTitle>
        </DialogHeader>
        <div class="py-4">
          <Input
            v-model="updatedMessageContent"
            placeholder="Enter updated message"
            class="w-full"
          />
        </div>
        <DialogFooter>
          <Button @click="closeUpdateModal" variant="outline">Cancel</Button>
          <Button @click="handleUpdateMessage" variant="default">Update</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="isDeleteModalOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Message</DialogTitle>
        </DialogHeader>
        <div class="py-4">
          <p>Are you sure you want to delete this message?</p>
          <div class="flex items-center mt-2">
            <input
              id="deleteForBoth"
              v-model="deleteForBoth"
              type="checkbox"
              class="mr-2"
            />
            <label for="deleteForBoth">Delete for both users</label>
          </div>
        </div>
        <DialogFooter>
          <Button @click="closeDeleteModal" variant="outline">Cancel</Button>
          <Button @click="handleDeleteMessage" variant="destructive">Delete</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- User Profile Modal -->
    <Sheet v-model:open="isProfileOpen">
      <SheetContent>
        <SheetHeader>
          <SheetTitle>User Profile</SheetTitle>
        </SheetHeader>
        <UserProfile 
          v-if="selectedUserProfile" 
          :user="selectedUserProfile" 
          :is-current-user="false"
        />
      </SheetContent>
    </Sheet>
  </div>
</template>