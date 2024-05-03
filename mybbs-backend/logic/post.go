package logic

import (
	"context"
	"go.uber.org/zap"
	"mybbs-backend/common"
	"mybbs-backend/config"
	"mybbs-backend/dto"
	"mybbs-backend/model"
	"strconv"
)

func CreatePost(userId int64, dto dto.PostDTO) error {
	db := config.GetDB()
	post := model.Post{
		UserID:  userId,
		Title:   dto.Title,
		Content: dto.Content,
	}
	if err := db.Create(&post).Error; err != nil {
		zap.L().Error("failed to create post", zap.Error(err))
		return err
	}
	return nil
}

func UpdatePost(postId int64, dto dto.UpdatePostDTO) error {
	db := config.GetDB()
	var post model.Post
	if err := db.Where("post_id = ?", postId).First(&post).Error; err != nil {
		return err
	}

	post.Title = dto.Title
	post.Content = dto.Content
	if err := db.Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func ReadPost(postId int64) (model.Post, error) {
	db := config.GetDB()
	var post model.Post
	if err := db.Where("post_id = ?", postId).Preload("User").Find(&post).Error; err != nil {
		return post, err
	}

	// 这里添加逻辑来增加浏览量：更新Redis和SQL数据库
	config.RDB.ZIncrBy(context.Background(), common.KeyPostScoreZSet, 1, strconv.FormatInt(postId, 10))
	post.ClickCount += 1
	db.Save(&post)

	return post, nil
}

func DeletePost(postId int64) error {
	db := config.GetDB()
	if err := db.Where("post_id = ?", postId).Delete(&model.Post{}).Error; err != nil {
		return err
	}
	return nil
}

func GetTopPostList(page int, pageSize int, order string) ([]model.Post, error) {
	db := config.GetDB()
	var posts []model.Post

	result := db.Preload("User").
		Scopes(config.Paginate(page, pageSize)).
		Order("click_count DESC").
		Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

func GetPostDetail(pid int32) model.Post {
	db := config.GetDB()
	var post model.Post
	db.Where(pid).Find(&post)
	return post
}

//func GetPostList(page int, pageSize int, order string) []model.Post {
//	if order == "ClickCount" {
//		return GetTopPostList(page, pageSize, order)
//	}
//	db := config.GetDB()
//	posts := make([]model.Post, 0)
//	db.Preload("User").
//		Scopes(config.Paginate(page, pageSize)).
//		Order(order + " DESC").
//		Find(&posts)
//
//	ids := make([]string, 0)
//	for i := range posts {
//		ids = append(ids, strconv.Itoa(int(posts[i].ID)))
//	}
//	ups, downs := GetVotes(ids)
//	for i := range posts {
//		posts[i].Up = ups[i]
//		posts[i].Down = downs[i]
//	}
//	return posts
//}

//// 按照帖子的浏览量获取最高浏览量的帖子
//func GetTopPostList(page int, pageSize int, order string) []model.Post {
//	rdb := config.RDB
//	start := (page - 1) * pageSize
//	end := start + pageSize - 1
//	// 查询出分数为前k个的post_id
//	ids, err := rdb.ZRevRange(context.Background(), common.KeyPostScoreZSet, int64(start), int64(end)).Result()
//	if err != nil {
//		zap.L().Error("GetTopPostList " + err.Error())
//	}
//	posts := getPostListIn(ids)
//	ups, downs := GetVotes(ids)
//	// TODO go不能使用for循环的post给结构体赋值 只能通过索引
//	//for i, post := range posts {
//	for i, _ := range posts {
//		posts[i].Up = ups[i]
//		posts[i].Down = downs[i]
//	}
//	return posts
//}

//
//// 统计帖子的点赞数和点踩数
//func GetVotes(ids []string) ([]int, []int) {
//	rdb := config.RDB
//	ctx := context.Background()
//	ups := make([]int, 0)
//	downs := make([]int, 0)
//	// TODO 可优化 for循环请求redis
//	for _, id := range ids {
//		result, _ := rdb.HVals(ctx, common.KeyPostVotedPrefix+id).Result()
//		var upTotal, downTotal int
//		for _, vote := range result {
//			if vote == "1" {
//				upTotal++
//			} else if vote == "-1" {
//				downTotal++
//			}
//		}
//		ups = append(ups, upTotal)
//		downs = append(downs, downTotal)
//	}
//	return ups, downs
//}

//func getPostListIn(ids []string) []model.Post {
//	db := config.DB
//	posts := make([]model.Post, 0)
//	// IN查询要指定ids顺序
//	db.Where("id IN ?", ids).Clauses(clause.OrderBy{
//		Expression: clause.Expr{SQL: "FIELD(id,?)",
//			Vars: []interface{}{ids}, WithoutParentheses: true}}).Find(&posts)
//	return posts
//}

//func CreatePost(userId int64, dto dto.PostDTO) {
//	db := config.DB
//	ctx := context.Background()
//	pipeline := config.RDB.TxPipeline()
//	post := model.Post{
//		UserID:  userId,
//		Title:   dto.Title,
//		Content: dto.Content,
//	}
//	// 入库
//	db.Create(&post)
//	// redis记录帖子发布时间
//	pipeline.ZAdd(ctx, common.KeyPostTimeZSet, &redis.Z{
//		Score:  float64(post.CreateTime.Unix()),
//		Member: post.ID,
//	})
//	// 设置帖子初始分值为帖子发布时间
//	pipeline.ZAdd(ctx, common.KeyPostScoreZSet, &redis.Z{
//		Score:  float64(post.CreateTime.Unix()),
//		Member: post.ID,
//	})
//	if cmds, err := pipeline.Exec(ctx); err != nil {
//		zap.L().Error("CreatePost redis操作异常 "+err.Error(), zap.Any("cmds", cmds))
//		return
//	}
//	zap.L().Info("redis记录发帖时间 " + fmt.Sprintf("%d", post.CreateTime.Unix()))
//}

//const (
//	OneWeekTime = 7 * 24 * 3600
//	Score       = 432
//)

// PostVoting 一天86400s 定义 200张投票=一天时间的分数
// 分数计算：帖子创建时间+点赞数*200
// TODO 如何做持久化？
//func PostVoting(userId int64, vote dto.VoteDTO) {
//	userIdStr := strconv.FormatInt(userId, 10)
//	postIdStr := strconv.FormatInt(int64(vote.PostID), 10)
//
//	rdb := config.RDB
//	pipeline := rdb.TxPipeline()
//	ctx := context.Background()
//	// 获取当前帖子的发表时间
//	//postTime := rdb.ZScore(ctx, common.KeyPostTimeZSet, postIdStr).Val()
//	//if float64(time.Now().Unix())-postTime > OneWeekTime {
//	//	panic(common.NewCustomError(common.CodeVoteTimeExpired))
//	//}
//	// 1. 查询当前用户是否有投票记录
//	t := rdb.HGet(ctx, common.KeyPostVotedPrefix+postIdStr, userIdStr).Val()
//	record, _ := strconv.Atoi(t) // 不存在key "" -> 0 会有err
//	// 2. 计算差值
//	diff := vote.Type - record
//	// 3. 更新当前帖子的分值
//	pipeline.ZIncrBy(ctx, common.KeyPostScoreZSet, float64(diff*Score), postIdStr)
//	// 4. 更新用户的投票记录
//	if vote.Type == 0 {
//		pipeline.HDel(ctx, common.KeyPostVotedPrefix+postIdStr, userIdStr)
//	} else {
//		pipeline.HSet(ctx, common.KeyPostVotedPrefix+postIdStr, userIdStr, strconv.Itoa(vote.Type))
//	}
//	if cmds, err := pipeline.Exec(ctx); err != nil {
//		zap.L().Error("PostVoting redis操作异常 "+err.Error(), zap.Any("cmds", cmds))
//	}
//}
