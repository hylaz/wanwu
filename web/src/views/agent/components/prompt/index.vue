<template>
  <div class="prompt-template-container">
    <div class="prompt-tabs">
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'builtIn' }"
        @click="activeTab = 'builtIn'"
      >
        {{ $t('agent.promptTemplate.builtIn') }}
      </div>
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'custom' }"
        @click="activeTab = 'custom'"
      >
        {{ $t('agent.promptTemplate.custom')}}
      </div>
    </div>

    <div class="cards-wrapper">
      <div v-if="showEmptyState" class="empty-state">
        <div class="empty-icon">
          <svg width="64" height="64" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 24H44" stroke="#D9D9D9" stroke-width="2" stroke-linecap="round"/>
            <path d="M20 32H44" stroke="#D9D9D9" stroke-width="2" stroke-linecap="round"/>
            <path d="M20 40H44" stroke="#D9D9D9" stroke-width="2" stroke-linecap="round"/>
            <path d="M16 12H48C49.1046 12 50 12.8954 50 14V50C50 51.1046 49.1046 52 48 52H16C14.8954 52 14 51.1046 14 50V14C14 12.8954 14.8954 12 16 12Z" stroke="#D9D9D9" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M14 18H50" stroke="#D9D9D9" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div class="empty-text">{{$t('agent.promptTemplate.emptyText')}}</div>
      </div>
      
      <div 
        v-else
        class="cards-container" 
        ref="cardsContainer"
        @scroll="handleScroll"
      >
        <div 
          class="scroll-button left" 
          v-if="showLeftButton"
          @click="scrollLeft"
        >
          <i class="el-icon-arrow-left"></i>
        </div>
        <div 
          v-for="(card, index) in currentCards" 
          :key="index"
          class="prompt-card"
          :ref="`card-${index}`"
          @click="handleCardClick(card)"
          @mouseenter="handleCardHover(card, index, $event)"
          @mouseleave="handleCardLeave"
        >
          <div class="card-title">{{ card.name }}</div>
          <div class="card-description">{{ card.prompt }}</div>
        </div>
        
        <div 
          v-if="!showEmptyState"
          class="prompt-card all-card"
          @click="handleAllClick"
        >
          <div class="all-card-content">
            <div class="all-card-text">{{$t('agent.promptTemplate.all')}}</div>
          </div>
        </div>
        
        <div 
          class="scroll-button right" 
          v-if="showRightButton"
          @click="scrollRight"
        >
          <i class="el-icon-arrow-right"></i>
        </div>
      </div>
      
      <div 
        v-if="hoveredCard"
        class="card-detail-panel"
        :class="{ 'panel-below': panelPositionBelow }"
        ref="detailPanel"
        :style="panelStyle"
        @mouseenter="handlePanelEnter"
        @mouseleave="handlePanelLeave"
      >
        <div class="detail-panel-title">{{ hoveredCard.name }}</div>
        <div class="detail-panel-content">
          <div class="detail-content-text" v-html="formatContent(hoveredCard.prompt)"></div>
        </div>
        <div class="detail-panel-footer">
          <el-button 
            type="primary" 
            class="insert-btn"
            @click="handleInsertPrompt(hoveredCard)"
            :disabled="!hoveredCard.prompt"
          >
            {{$t('agent.promptTemplate.insertPrompt')}}
          </el-button>
        </div>
      </div>
    </div>
    
    <promptDialog ref="promptDialog"  @tabChange="handleTabChange"/>
  </div>
</template>

<script>
import promptDialog from './promptDialog.vue';
import { md } from "@/mixins/marksown-it.js";
import {getPromptTemplateList,getPromptBuiltInList} from "@/api/promptTemplate"
export default {
  name: 'PromptTemplate',
  inject: ['getPrompt'],
  components: {
    promptDialog
  },
  data() {
    return {
      activeTab: 'builtIn',
      showLeftButton: false,
      showRightButton: true,
      hoveredCard: null,
      hoveredCardIndex: null,
      hoverTimer: null,
      leaveTimer: null,
      isMouseOnPanel: false,
      isDefaultShow: false,
      panelPositionBelow: false,
      panelLeft: '50%',
      panelTop: '-420px',
      recommendedCards: [
       
      ],
      personalCards: [
      ]
    }
  },
  computed: {
    currentCards() {
      const cards = this.activeTab === 'builtIn' ? this.recommendedCards : this.personalCards;
      return cards.slice(0, 6);
    },
    showEmptyState() {
      if (this.activeTab === 'builtIn') {
        return !this.recommendedCards || this.recommendedCards.length === 0;
      }
      return !this.personalCards || this.personalCards.length === 0;
    },
    panelStyle() {
      if (!this.hoveredCard) return {};
      return {
        position: 'fixed',
        left: this.panelLeft,
        top: this.panelTop,
        transform: 'translateX(-50%)',
        zIndex: 1000
      };
    },
    cardContent() {
      if (!this.hoveredCard) return '';
      if (this.hoveredCard.content) {
        return this.hoveredCard.content;
      }
      const promptDialog = this.$refs.promptDialog;
      if (promptDialog && promptDialog.templateList) {
        const template = promptDialog.templateList.find(t => 
          t.name === this.hoveredCard.title || 
          t.name === this.hoveredCard.name
        );
        if (template && template.content) {
          return template.content;
        }
      }
      
      return this.hoveredCard.description || '';
    }
  },
  created(){
    this.getPromptBuiltInList()
    this.getPromptTemplateList()
  },
  mounted() {
    this.checkScrollButton();
    window.addEventListener('resize', this.checkScrollButton);
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.checkScrollButton);
    if (this.hoverTimer) {
      clearTimeout(this.hoverTimer);
    }
    if (this.leaveTimer) {
      clearTimeout(this.leaveTimer);
    }
  },
  methods: {
    handleTabChange(tab){
      const cards = tab === 'builtIn' ? this.recommendedCards : this.personalCards;
      this.$refs.promptDialog.showDiglog(cards, tab);
    },
    getPromptTemplateList(name=''){
      getPromptTemplateList({name}).then(res =>{
        if(res.code === 0){
          this.personalCards = res.data.list || []
        }
      })
    },
    getPromptBuiltInList(name=''){
      getPromptBuiltInList({name,category:'all'}).then(res =>{
        if(res.code === 0){
            this.recommendedCards = res.data.list || []
            this.$nextTick(() => {
              this.checkScrollButton();
            });
        }
      })
    },
    handleCardClick(card) {
      this.$emit('card-click', card);
    },
    handleCardHover(card, index, event) {
      if (this.leaveTimer) {
        clearTimeout(this.leaveTimer);
        this.leaveTimer = null;
      }
      if (this.hoverTimer) {
        clearTimeout(this.hoverTimer);
      }
      this.hoverTimer = setTimeout(() => {
        this.hoveredCard = card;
        this.hoveredCardIndex = index;
        this.isMouseOnPanel = false;
        this.isDefaultShow = false;
        this.updatePanelPosition(event);
      }, 200);
    },
    handleCardLeave() {
      if (this.hoverTimer) {
        clearTimeout(this.hoverTimer);
        this.hoverTimer = null;
      }
      this.leaveTimer = setTimeout(() => {
        if (!this.isMouseOnPanel && !this.isDefaultShow) {
          this.hoveredCard = null;
          this.hoveredCardIndex = null;
        }
      }, 150);
    },
    handlePanelEnter() {
      this.isMouseOnPanel = true;
      if (this.hoverTimer) {
        clearTimeout(this.hoverTimer);
        this.hoverTimer = null;
      }
      if (this.leaveTimer) {
        clearTimeout(this.leaveTimer);
        this.leaveTimer = null;
      }
    },
    handlePanelLeave() {
      this.isMouseOnPanel = false;
      if (!this.isDefaultShow) {
        this.hoveredCard = null;
        this.hoveredCardIndex = null;
      }
    },
    updatePanelPosition(event) {
      this.$nextTick(() => {
        if (this.hoveredCardIndex === null || this.hoveredCardIndex === undefined) return;
        
        const cardRef = this.$refs[`card-${this.hoveredCardIndex}`];
        if (!cardRef || !cardRef[0]) return;
        
        const cardRect = cardRef[0].getBoundingClientRect();
        const panelWidth = 300;
        const panelHeight = 240;
        
        const cardCenterX = cardRect.left + cardRect.width / 2;
        const cardTop = cardRect.top;
        const cardBottom = cardRect.bottom;
        
        let panelTop = cardTop - panelHeight - 10;
        let isBelow = false;
        
        if (panelTop < 10) {
          panelTop = cardBottom + 10;
          isBelow = true;
        }
        
        this.panelPositionBelow = isBelow;
        
        let panelLeft = cardCenterX;
        if (panelLeft - panelWidth / 2 < 20) {
          panelLeft = panelWidth / 2 + 20;
        } else if (panelLeft + panelWidth / 2 > window.innerWidth - 20) {
          panelLeft = window.innerWidth - panelWidth / 2 - 20;
        }
        
        this.panelLeft = `${panelLeft}px`;
        this.panelTop = `${panelTop}px`;
      });
    },
    formatContent(content) {
      if (!content) return '';
      if (content.includes('##') || content.includes('#') || content.includes('**') || content.includes('*')) {
        return md.render(content);
      }
      return `<p>${content.replace(/\n/g, '<br/>')}</p>`;
    },
    handleInsertPrompt(card) {
      const content = this.getCardContent(card);
      if (content) {
        this.getPrompt(card.prompt)
        this.hoveredCard = null;
        this.hoveredCardIndex = null;
        this.isDefaultShow = false;
        this.$message.success(this.$t('agent.promptTemplate.insertSuccess'));
      }
    },
    getCardContent(card) {
      if (!card) return null;
      
      if (card.prompt) {
        return card.prompt;
      }
      
      const promptDialog = this.$refs.promptDialog;
      if (promptDialog && promptDialog.templateList) {
        const template = promptDialog.templateList.find(t => t.name === card.name);
        if (template && template.prompt) {
          return template.prompt;
        }
      }
      return null;
    },
    handleAllClick() {
      const cards = this.activeTab === 'builtIn' ? this.recommendedCards : this.personalCards;
      this.$refs.promptDialog.showDiglog(cards,this.activeTab);
    },
    scrollLeft() {
      const container = this.$refs.cardsContainer;
      if (container) {
        const scrollAmount = 300;
        container.scrollBy({
          left: -scrollAmount,
          behavior: 'smooth'
        });
      }
    },
    scrollRight() {
      const container = this.$refs.cardsContainer;
      if (container) {
        const scrollAmount = 300;
        container.scrollBy({
          left: scrollAmount,
          behavior: 'smooth'
        });
      }
    },
    handleScroll() {
      this.checkScrollButton();
    },
    checkScrollButton() {
      this.$nextTick(() => {
        const container = this.$refs.cardsContainer;
        if (!container) return;
        const originalCards = this.activeTab === 'builtIn' ? this.recommendedCards : this.personalCards;
        const hasMoreCards = originalCards && originalCards.length > 6;
        const canScroll = container.scrollWidth > container.clientWidth;
        const scrollLeft = container.scrollLeft;
        const scrollWidth = container.scrollWidth;
        const clientWidth = container.clientWidth;
        
        const isAtStart = scrollLeft <= 10;
        const isAtEnd = scrollLeft + clientWidth >= scrollWidth - 10;
        
        if (!canScroll && !hasMoreCards) {
          this.showLeftButton = false;
          this.showRightButton = false;
        } else {
          if (isAtStart) {
            this.showLeftButton = false;
            this.showRightButton = canScroll || hasMoreCards;
          } 
          else if (isAtEnd) {
            this.showLeftButton = canScroll;
            this.showRightButton = false;
          } 
          else {
            this.showLeftButton = true;
            this.showRightButton = true;
          }
        }
      });
    }
  },
  watch: {
    activeTab() {
      this.$nextTick(() => {
        this.checkScrollButton();
      });
    },
    recommendedCards() {
      this.$nextTick(() => {
        this.checkScrollButton();
      });
    },
    personalCards() {
      this.$nextTick(() => {
        this.checkScrollButton();
      });
    }
  }
}
</script>

<style lang="scss" scoped>
.prompt-template-container {
  position:absolute;
  bottom:0;
  left:0;
  right:0;
  width: 100%;
  padding: 10px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.prompt-tabs {
  display: flex;
  gap: 10;
  .tab-item {
    padding: 3px 8px;
    cursor: pointer;
    color: #303133;
    font-size: 14px;
    transition: all 0.3s;
    border: none;
    white-space: nowrap;
    border-radius: 4px 4px 0 0;
    
    &:hover {
      color: $color;
    }
    
    &.active {
      color: $color;
      background: #E0E7FF;
      font-weight: 500;
    }
  }
}

.cards-wrapper {
  position: relative;
  flex: 1;
  overflow: visible;
  width: 100%;
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 10px 0;
  
  .empty-icon {
    margin-bottom: 16px;
    opacity: 0.6;
    
    svg {
      display: block;
    }
  }
  
  .empty-text {
    font-size: 14px;
    color: #606266;
    text-align: center;
    line-height: 1.5;
  }
}

.cards-container {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  overflow-y: visible;
  padding: 10px 0;
  scroll-behavior: smooth;
  position: relative;
  align-items: stretch;
  width: 100%;
  max-width: 100%;
  
  &::-webkit-scrollbar {
    display: none;
  }
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.card-detail-panel {
  width: 300px;
  height: 240px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  animation: fadeInUp 0.3s ease-out;
  padding: 10px;
  position: relative;
  &::after {
    content: '';
    position: absolute;
    bottom: -8px;
    left: 50%;
    transform: translateX(-50%);
    width: 0;
    height: 0;
    border-left: 8px solid transparent;
    border-right: 8px solid transparent;
    border-top: 8px solid #fff;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
  }
  
  &.panel-below {
    &::after {
      bottom: auto;
      top: -8px;
      border-top: none;
      border-bottom: 8px solid #fff;
      filter: drop-shadow(0 -2px 4px rgba(0, 0, 0, 0.1));
    }
  }
  
  .detail-panel-title {
    padding: 5px 15px;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    flex-shrink: 0;
  }
  
  .detail-panel-content {
    flex: 1;
    padding: 10px 15px;
    overflow-y: auto;
    overflow-x: hidden;
    min-height: 0;
    width: 100%;
    box-sizing: border-box;
    
    .detail-content-text {
      margin: 0;
      font-size: 14px;
      line-height: 1.8;
      color: #303133;
      word-wrap: break-word;
      overflow-wrap: break-word;
      word-break: break-word;
      white-space: normal;
      max-width: 100%;
      box-sizing: border-box;
      
      /deep/ {
        * {
          max-width: 100%;
          word-wrap: break-word;
          overflow-wrap: break-word;
          word-break: break-word;
        }
        
        strong {
          background-color: #F2F0FF;
          padding: 2px 6px;
          border-radius: 4px;
          color: #7C3AED;
          font-weight: 500;
          display: inline;
          word-break: break-word;
        }
        
        h1, h2 {
          margin-top: 16px;
          margin-bottom: 8px;
          font-weight: 600;
          word-break: break-word;
          overflow-wrap: break-word;
          
          &:first-child {
            margin-top: 0;
          }
        }
        
        ul, ol {
          margin: 8px 0;
          padding-left: 24px;
          word-break: break-word;
          
          li {
            margin: 4px 0;
            word-break: break-word;
            overflow-wrap: break-word;
          }
        }
        
        p {
          margin: 8px 0;
          word-break: break-word;
          overflow-wrap: break-word;
        }
        
        code, pre {
          word-break: break-all;
          white-space: pre-wrap;
          max-width: 100%;
        }
      }
    }
  }
  
  .detail-panel-footer {
    height: 50px;
    padding: 10px;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-shrink: 0;
    .insert-btn {
      width: 100%;
      border: none;
      color: $color;
      font-weight: 500;
      padding: 10px 24px;
      border-radius: 6px;
      transition: all 0.3s;
      &:active {
        transform: translateY(0);
      }
    }
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translate(-50%, 10px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}

.prompt-card {
  flex-shrink: 0;
  width: 200px;
  background: #fff;
  border-radius: 8px;
  padding: 10px;
  box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
  
  &:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border-color: $color;
    transform: translateY(-2px);
  }
  
  &.all-card {
    background: #fff;
    border: 1px solid transparent;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    
    &:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      border-color: $color;
      transform: translateY(-2px);
    }
    
    .all-card-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: 100%;
      padding: 20px;
      
      .all-card-text {
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
    }
  }
  
  .card-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 10px;
    line-height: 1.2;
  }
  
  .card-description {
    font-size: 13px;
    color: #606266;
    line-height: 1.6;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-word;
    max-height: calc(1.6em * 3);
    min-height: calc(1.6em * 3);
  }
}

.scroll-button {
  position: sticky;
  top: auto;
  align-self: center;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  min-width: 32px;
  min-height: 32px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 5;
  transition: all 0.3s;
  border: 1px solid #e4e7ed;
  
  &.left {
    left: 0;
    margin-right: 8px;
  }
  
  &.right {
    right: 0;
    margin-left: 8px;
  }
  
  &:hover {
    background: #f5f7fa;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
  
  i {
    font-size: 14px;
    color: #606266;
    font-weight: bold;
    line-height: 1;
  }
  
  &:hover i {
    color: $color;
  }
}

@media (max-width: 768px) {
  .prompt-card {
    min-width: 240px;
    max-width: 240px;
  }
}
</style>
