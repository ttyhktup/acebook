import "./Post.css";
import image from "/src/static/img/cross.png";
import like from "/src/static/img/heart.png"
import { useState, useEffect } from "react";
import React from "react";
import Comment from "../Comment/Comment";
import { createComments, getComments, deleteComments } from "../../services/comments"

const Post = ({ post, onLike, user, onDelete, token, loggedUserID }) => {
  const [comments, setComments] = useState([]);
  const [comment, setComment] = useState("");
  const [isClicked, setIsClicked] = useState(false);
  const [commentUserID, setCommentUserID] = useState("");
  const { Fragment } = React;

  const handleLikeClick = () => {
    onLike(post._id);
    setIsClicked(true)
    setTimeout(() => {
      setIsClicked(false); // Revert the state back to false
    }, 100);
  };

  const handleDeleteClick = () => {
    onDelete(post._id);
  };

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      getComments(post._id, token)
          .then((data) => {
            const sortedPosts = data.comments.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
            setComments(sortedPosts)
            localStorage.setItem("token", data.token);
            setCommentUserID(data.loggedUserID);
          })
          .catch((err) => {
            console.error(err);
          });
    }
  }, [post._id]);


  const handleSubmitComment = async (event) => {
    event.preventDefault();
    try {
      await createComments(token, comment, post._id);
      const updatedComments = await getComments(post._id, token);
      const sortedComments = updatedComments.comments.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      setComments(sortedComments);
      setComment("");
      localStorage.setItem("token", updatedComments.token);
    } catch (err) {
      console.error(err);
    }
  };

  const handleDeleteComment = async (commentId) => {
    try {
      await deleteComments(token, post._id, commentId);
      const updatedComments = await getComments(post._id, token);
      const sortedComments = updatedComments.comments.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      setComments(sortedComments);
    } catch (err) {
      console.error(err);
    }
  };

  const handleCommentChange = (event) => {
    setComment(event.target.value);
  };

  const buttonClass = isClicked ? 'like-button clicked' : 'like-button'

  return (
      <div className="posts-and-comments" key={post._id}>
        <div className="post">
          <div className="post-user">
            <div className="post-user-info">
              <img className="user-image" src={user.image} alt="image" />
              <p>{user.username}</p>
            </div>
            {loggedUserID == user.user_id && <img className="delete-button" src={image} alt="delete" onClick={handleDeleteClick} />}
          </div>
          <div className="post-content">
            <div className="post-message"><p>{post.message}</p></div>
            <div className="likes">
              <img className={buttonClass} src={like} onClick={handleLikeClick}/>
              <p>Likes: {post.likes}</p>
            </div>
          </div>
        </div>
            
        <div className="comments">
          <h4>Comments</h4>
          <form className="create-comment" onSubmit={handleSubmitComment}>
              <input
                  className="comment-input"
                  type="text"
                  onChange={handleCommentChange}
                  value={comment}
                  placeholder="Add a comment..."
              />
              <button className="comment-submit" type="submit">Submit</button>
          </form>
          
          {comments.length < 1 && <div className="no-comments"> No comments yet </div>}
          {comments
              .map((comment) => (
                <Fragment key={comment._id}>
                    <Comment commentUserID={commentUserID} comment={comment} onDelete={handleDeleteComment} />
                  </Fragment>
              ))}
        </div>

      </div>
  );
};

export default Post;