<template>
  <div class="audio-recorder">
    <Button
      v-if="!isRecording"
      @click="startRecording"
      variant="outline"
      size="icon"
      type="button"
      :disabled="!isSupported"
    >
      <MicIcon class="h-5 w-5" />
      <span class="sr-only">Start Recording</span>
    </Button>
    <div v-else class="flex items-center space-x-2">
      <Button
        @click="stopRecording"
        variant="destructive"
        size="icon"
        type="button"
      >
        <StopCircleIcon class="h-5 w-5" />
        <span class="sr-only">Stop Recording</span>
      </Button>
      <span class="text-sm font-medium">{{ formatTime(recordingTime) }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue';
import { MicIcon, StopCircleIcon } from 'lucide-vue-next';
import { Button } from '@/components/ui/button';

const isSupported = ref(!!navigator.mediaDevices && !!navigator.mediaDevices.getUserMedia);
const isRecording = ref(false);
const mediaRecorder = ref<MediaRecorder | null>(null);
const audioChunks = ref<Blob[]>([]);
const recordingTime = ref(0);
const recordingInterval = ref<number | null>(null);

const emit = defineEmits<{
  (e: 'sendAudio', blob: Blob): void
}>();

const startRecording = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    mediaRecorder.value = new MediaRecorder(stream);
    
    mediaRecorder.value.ondataavailable = (event) => {
      audioChunks.value.push(event.data);
    };
    
    mediaRecorder.value.onstop = () => {
      const audioBlob = new Blob(audioChunks.value, { type: 'audio/wav' });
      emit('sendAudio', audioBlob);
      audioChunks.value = [];
    };
    
    mediaRecorder.value.start();
    isRecording.value = true;
    recordingTime.value = 0;
    recordingInterval.value = window.setInterval(() => {
      recordingTime.value++;
    }, 1000);
  } catch (error) {
    console.error('Error accessing microphone:', error);
  }
};

const stopRecording = () => {
  if (mediaRecorder.value && isRecording.value) {
    mediaRecorder.value.stop();
    isRecording.value = false;
    clearInterval(recordingInterval.value!);
    mediaRecorder.value.stream.getTracks().forEach(track => track.stop());
  }
};

const formatTime = (seconds: number): string => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
};

onUnmounted(() => {
  if (isRecording.value) {
    stopRecording();
  }
});
</script>